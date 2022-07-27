-- Create syntax for 'attendees'

CREATE TABLE `attendees` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` tinytext NOT NULL,
  `email` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `tour` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `checked_in` int NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;