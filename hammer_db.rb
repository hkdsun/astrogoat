require 'mysql2'

conn = Mysql2::Client.new(host: "127.0.0.1", username: "root", port: 21001)

conn.query("CREATE DATABASE IF NOT EXISTS test")
conn.query("CREATE TABLE IF NOT EXISTS test.test (id INT, PRIMARY KEY(id))")

loop do
  conn.query("INSERT IGNORE INTO test.test VALUES(#{Integer(rand * 99999999999)})")
end
