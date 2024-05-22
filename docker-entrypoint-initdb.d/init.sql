create table users(
                      user_id SERIAL PRIMARY KEY,
                      username VARCHAR(50) UNIQUE NOT NULL,
                      password_hash TEXT NOT NULL,
                      role VARCHAR(10) NOT NULL CHECK (role IN ('loader', 'customer'))
);

create table customers(
                          customer_id INT UNIQUE NOT NULL,
                          starting_capital DECIMAL(10, 2) CHECK (starting_capital BETWEEN 10000 AND 100000),
                          current_capital DECIMAL(10, 2) CHECK (starting_capital BETWEEN 10000 AND 100000),
                          FOREIGN KEY (customer_id) REFERENCES Users(user_id)
);

create table loaders(
                        loader_id INT UNIQUE NOT NULL,
                        max_weight DECIMAL(5, 2) CHECK (max_weight BETWEEN 5 AND 30),
                        drunkenness BOOLEAN NOT NULL,
                        fatigue DECIMAL(5, 2) CHECK (fatigue BETWEEN 0 AND 100) DEFAULT 0,
                        salary DECIMAL(10, 2) CHECK (salary BETWEEN 10000 AND 30000),
                        FOREIGN KEY (loader_id) REFERENCES Users(user_id)
);

create table tasks(
                      task_id SERIAL PRIMARY KEY,
                      customer_id INT NOT NULL,
                      task_name VARCHAR(100) NOT NULL,
                      weight DECIMAL(5, 2) CHECK (weight BETWEEN 10 AND 80),
                      completed BOOLEAN DEFAULT FALSE,
                      FOREIGN KEY (customer_id) REFERENCES Users(user_id)

);


create table loaders_tasks(
                              loader_id INT NOT NULL,
                              task_id INT NOT NULL,
                              FOREIGN KEY (loader_id) REFERENCES Users(user_id),
                              FOREIGN KEY (task_id) REFERENCES Tasks(task_id),
                              UNIQUE (loader_id, task_id)
);