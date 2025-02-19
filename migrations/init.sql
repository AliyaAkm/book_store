CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(100) NOT NULL UNIQUE,
                       password VARCHAR(100) NOT NULL
);

INSERT INTO users (username, password) VALUES ('Aliya', '123');

// думаю тут неправильная реализация. несмотря на этот код я создавала в самом бд таблицу.
// таблицу book я не создавала, но в бд она правильно отображается
//  аоаоаоаоао