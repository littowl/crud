-- up.sql
CREATE TABLE IF NOT EXISTS users_articles (
    Author_Id INT REFERENCES users,
    Article_Id INT REFERENCES articles,
    PRIMARY KEY (Author_Id, Article_Id)
)