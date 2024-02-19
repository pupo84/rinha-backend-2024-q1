package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"rinha-backend/api/dto"

	_ "github.com/lib/pq"
)

type Database struct {
	Connection *sql.DB
}

func NewDatabase() *Database {
	host := "localhost"
	port := 5432
	username := "postgres"
	password := "postgres"
	dbname := "rinha"
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname)
	db, err := sql.Open("postgres", conn)
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)
	if err != nil {
		log.Fatalf("Could not estabilish connection with postgres due to %v\n", err.Error())
	}
	return &Database{db}
}

func (db *Database) GetUser(ctx context.Context, userID int64) (user dto.User, err error) {
	err = db.Connection.QueryRow(`select id, name, email, document, "limit" from users where id = $1`, userID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Document, &user.Limit)
	if err != nil {
		fmt.Printf("Could not get user due to %v\n", err.Error())
		return
	}
	return
}

func (db *Database) CreateUser(ctx context.Context, request dto.CreateUserRequest) (user dto.User, err error) {
	var result sql.Result

	txn, err := db.Connection.Begin()
	if err != nil {
		log.Fatalf("Could not get transaction with postgres due to %v\n", err.Error())
	}

	insertStmt := "insert into users values (nextval('users_seq'), $1, $2, $3, $4, now())"
	result, err = txn.ExecContext(ctx, insertStmt, request.Name, request.Email, request.Document, request.Limit)
	if err != nil {
		fmt.Printf("Could not create user due to %s", err.Error())
		txn.Rollback()
		return user, err
	}

	rows, err := result.RowsAffected()
	if rows == 0 || err != nil {
		fmt.Printf("Could not create user due to %s", err.Error())
		return
	}

	queryStmt := `select id, name, email, document, "limit" from users where document = $1`
	if err = txn.QueryRow(queryStmt, request.Document).
		Scan(&user.ID, &user.Name, &user.Email, &user.Document, &user.Limit); err != nil {
		fmt.Printf("Could not get user due to %v\n", err.Error())
		txn.Rollback()
		return
	}

	insertBalanceStmt := "insert into balances values (nextval('balances_seq'), $1, $2, now())"
	result, err = txn.ExecContext(ctx, insertBalanceStmt, user.ID, 0)
	if err != nil {
		fmt.Printf("Could not insert balance for user due to %v\n", err.Error())
		txn.Rollback()
		return
	}

	rows, err = result.RowsAffected()
	if rows == 0 || err != nil {
		fmt.Printf("Could not insert balance for user due to %v\n", err.Error())
		return
	}

	txn.Commit()

	return
}

func (db *Database) MakeTransaction(ctx context.Context, userID int64, transaction dto.Transaction) bool {
	var (
		result sql.Result
		err    error
	)

	updateCreditStmt := "update balances set amount = amount + $1, updated_at = now()"
	updateDebitStmt := `update balances b 
	                       set amount = amount - $2, updated_at = now() 
						  from users u
						 where b.user_id = u.id
						   and b.user_id = $1
						   and b.amount + u.limit >= $2`

	txn, err := db.Connection.Begin()
	if err != nil {
		log.Fatalf("Could not get transaction with postgres due to %v\n", err.Error())
	}

	if transaction.Nature == "c" {
		result, err = txn.ExecContext(ctx, updateCreditStmt, transaction.Amount)
	} else if transaction.Nature == "d" {
		result, err = txn.ExecContext(ctx, updateDebitStmt, userID, transaction.Amount)
	} else {
		fmt.Printf("Invalid transaction type %s", transaction.Nature)
		txn.Rollback()
		return false
	}
	if err != nil {
		fmt.Printf("Could not update funds due to %v\n", err.Error())
		txn.Rollback()
		return false
	}

	rows, err := result.RowsAffected()
	if rows == 0 || err != nil {
		fmt.Println("Could not update funds")
		txn.Rollback()
		return false
	}

	insertTxnSmtm := "insert into transactions values (nextval('transactions_seq'), $1, $2, $3, $4, now())"
	result, err = txn.ExecContext(ctx, insertTxnSmtm, userID, transaction.Nature, transaction.Amount, transaction.Description)
	if err != nil {
		fmt.Printf("Could not update balance. Error creating transaction %v\n", err.Error())
		txn.Rollback()
		return false
	}

	rows, err = result.RowsAffected()
	if rows == 0 || err != nil {
		fmt.Printf("Could not insert transaction due to %v\n", err.Error())
		txn.Rollback()
		return false
	}

	txn.Commit()

	return true
}

func (db *Database) GetBalance(ctx context.Context, userId int64) (balance dto.Balance, err error) {
	err = db.Connection.QueryRow("select u.limit, b.amount, b.updated_at from users u inner join balances b on u.id = b.user_id where b.id = $1", userId).
		Scan(&balance.Limit, &balance.Balance, &balance.UpdatedAt)
	if err != nil {
		return
	}
	return
}

func (db *Database) GetStatement(ctx context.Context, userID int64) (statement dto.TransactionResponse, err error) {
	statementStmt := `select amount, "type", description, created_at from transactions where user_id = $1`
	rows, err := db.Connection.Query(statementStmt, userID)
	if err != nil {
		return
	}

	txns := make([]dto.Transaction, 0)
	for rows.Next() {
		txn := dto.Transaction{}
		rows.Scan(&txn.Amount, &txn.Nature, &txn.Description, &txn.CreatedAt)
		txns = append(txns, txn)
	}

	balance, _ := db.GetBalance(ctx, userID)
	balanceResponse := dto.BalanceResponse{
		Total:         balance.Balance,
		StatementDate: balance.UpdatedAt,
		Limit:         balance.Limit,
	}
	statement.BalanceResponse = balanceResponse
	statement.Transactions = txns

	return
}
