package sy

import (
	"context"
	"errors"
	"fmt"
	"grpcserv/pb"
	"log"
	"os"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	UserRepo struct {
		pool *pgxpool.Pool
	}
	User struct {
		ID         int64      `json:"id" db:"id"`
		Yus_code   string     `json:"yus_code" db:"yus_code"`
		Yus_passwd string     `json:"yus_passwd" db:"yus_passwd"`
		Yus_name   string     `json:"yus_name" db:"yus_name"`
		Yus_memo   string     `json:"yus_memo" db:"yus_memo"`
		Yus_active bool       `json:"yus_active" db:"yus_active"`
		Yus_system bool       `json:"yus_system" db:"yus_system"`
		CreatedAt  *time.Time `json:"create_at" db:"createdat"`
		UpdatedAt  *time.Time `json:"update_at" db:"updatedat"`
	}
)

func NewUserRepo(pool *pgxpool.Pool, ctx context.Context) *UserRepo {
	if pool == nil {
		var err error
		pool, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
			os.Exit(1)
		}
	}
	if err := pool.Ping(ctx); err != nil {
		log.Println("Unable to connect to database:", err)
		return nil
	}
	return &UserRepo{pool: pool}
}

func (r UserRepo) GetUserById(ctx context.Context, id int64) (User, error) {
	query := `SELECT id, yus_code, yus_passwd, yus_name, yus_memo, yus_active, yus_system, createdat, updatedat
		FROM yus_users WHERE id = $1`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		log.Println("Query error:", err)
		return User{}, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, nil
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r UserRepo) GetUserByCode(ctx context.Context, code string) (User, error) {
	query := `SELECT id, yus_code, yus_passwd, yus_name, yus_memo, yus_active, yus_system, createdat, updatedat
		FROM yus_users WHERE yus_code = $1`
	rows, err := r.pool.Query(ctx, query, code)
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, nil
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r UserRepo) GetUserByLogin(ctx context.Context, username, password string) (User, error) {
	query := `SELECT id, yus_code, yus_passwd, yus_name, yus_memo, yus_active, yus_system, createdat, updatedat
		FROM yus_users WHERE yus_code = $1 AND yus_passwd = $2`
	rows, err := r.pool.Query(ctx, query, username, password)
	if err != nil {
		log.Println("Query error:", err)
		return User{}, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, nil
	}
	if err != nil {
		return User{}, err
	}
	user.Yus_passwd = "" // Clear password before returning
	return user, nil
}

func (r UserRepo) AddUser(ctx context.Context, in *pb.UserAddRequest) (int64, error) {
	var id int
	query := `INSERT INTO yus_users (yus_code, yus_passwd, yus_name, yus_memo, yus_active)
		VALUES ($1, $2, $3, $4, true)
		RETURNING id`
	err := r.pool.QueryRow(
		ctx,
		query,
		in.GetYusCode(),
		in.GetYusPasswd(),
		in.GetYusName(),
		in.GetYusMemo(),
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (r UserRepo) SetUser(ctx context.Context, in *pb.UserSetRequest) error {
	query := `UPDATE yus_users SET yus_code = $1, yus_passwd = $2, yus_memo = $3, yus_active = $4
		WHERE id = $5 RETURNING id, yus_code, yus_passwd, yus_memo, yus_active, yus_system, createdat, updatedat`
	_, err := r.pool.Exec(
		ctx,
		query,
		in.GetYusCode(),
		in.GetYusPasswd(),
		in.GetYusMemo(),
		in.GetYusActive(),
		in.GetId(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepo) GetUserList(ctx context.Context, filter string, page int64, lpp int64) ([]User, error) {
	var rows pgx.Rows
	var err error
	if len(filter) == 0 {
		query := `SELECT id, yus_code, yus_passwd, yus_name, yus_memo, yus_active, yus_system, createdat, updatedat
			FROM yus_users ORDER BY id OFFSET $1 LIMIT $2`
		rows, err = r.pool.Query(ctx, query, page*lpp, lpp)
	} else {
		ftxt := fmt.Sprintf("%%%s%%", filter)
		query := `SELECT id, yus_code, yus_passwd, yus_name, yus_memo, yus_active, yus_system, createdat, updatedat
			FROM yus_users WHERE yus_name ilike $1 ORDER BY id OFFSET $2 LIMIT $3`
		rows, err = r.pool.Query(ctx, query, ftxt, page*lpp, lpp)
	}
	fmt.Println("Repo.GetUserList2: ", rows.RawValues(), rows.Err(), err)
	if err != nil {
		return []User{}, err
	}
	users, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[User],
	)
	fmt.Println("Repo.GetUserList3: ", users)

	return users, nil
}

func (r UserRepo) GetUserCount(ctx context.Context, filter string) (int64, error) {
	ftxt := fmt.Sprintf("%%%s%%", filter)
	query := `SELECT COUNT(*) FROM yus_users WHERE yus_name ilike $1`
	var count int64
	err := r.pool.QueryRow(ctx, query, ftxt).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r UserRepo) DeleteUserById(ctx context.Context, id int64) error {
	query := `DELETE FROM yus_users WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r UserRepo) DeleteUserByCode(ctx context.Context, code string) error {
	query := `DELETE FROM yus_users WHERE yus_code = $1`
	_, err := r.pool.Exec(ctx, query, code)
	if err != nil {
		return err
	}
	return nil
}
