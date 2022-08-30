CREATE TABLE `book` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `price` float DEFAULT NULL,
  `author` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `image_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


INSERT INTO `book` (`id`, `name`, `price`, `author`, `description`, `image_url`)
VALUES
	(1, "Book 1: Harry Potter and the Sorcerer's Stone", 530, 'JK rowling', "Harry Potter and the Sorcerer's Stone", 'http://www.adviceforyou.co.th/blog/wp-content/uploads/2011/12/harry-potter.jpeg'),
  (2, 'Book 2: Harry Potter and the Chamber of Secrets', 530, 'JK rowling', "Book 2: Harry Potter and the Chamber of Secrets", 'https://content.time.com/time/2007/harry_potter/hp_books/chamber_of_secrets.jpg');