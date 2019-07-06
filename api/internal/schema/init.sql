ALTER DATABASE example
    DEFAULT CHARACTER SET utf8
    DEFAULT COLLATE utf8_general_ci;

CREATE TABLE `users` (
    `id` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
    `email` VARCHAR(100) NOT NULL DEFAULT '',
    `password` VARCHAR(100) NOT NULL DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB CHARSET=utf8;

INSERT INTO users (email, password) VALUES
    ('test@test.com', '$2a$10$eNxD0bfWdeWw3o3gLGVjNeiw/H0/KVaz6wh/UkFmKHm2ZJXOhvvVW');
