#!/usr/bin/env gorun

package src

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func newConnection() *sql.DB {
  database, err := sql.Open("sqlite3", "./test.db")
  if err != nil {
    panic(err)
  }
  statement, _ := database.Prepare(
      `CREATE TABLE IF NOT EXITS todos (
      id  INTEGER PRIMARY KEY AUTOINCREMENT,
      name  TEXT,
      completed  BOOLEAN,
      createdAt  DATETIME
    )`)

  statement.Exec()
  return database
}
