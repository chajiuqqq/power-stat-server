CREATE TABLE `report` (
                               `id` int NOT NULL AUTO_INCREMENT,
                               `time` datetime  not null,
                               `total_start_time` datetime  not null,
                               `total` float(10,3)  not null,
  `yesterday` float(10,3)  not null,
  `today` float(10,3)  not null,
  `period` int  not null,
  `power` int  not null,
  `apparent_power` int  not null,
  `reactive_power` int  not null,
  `factor` float(4,2)  not null,
  `frequency` int  not null,
  `voltage` int  not null,
  `current` float(10,3)  not null,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
