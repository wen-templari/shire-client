CREATE TABLE `group_user`  (
  `groupId` int NOT NULL,
  `userId` int NOT NULL,
  PRIMARY KEY (`groupId`, `userId`)
);

CREATE TABLE `message`  (
  `from` int NOT NULL DEFAULT 0,
  `to` int NOT NULL DEFAULT 0,
  `content` varchar(255) NOT NULL DEFAULT "",
  `groupId` int NOT NULL DEFAULT 0,
  `time` datetime NULL
);

CREATE TABLE `user`  (
  `id` int NOT NULL ,
  `name` varchar(255) NOT NULL DEFAULT "",
  `address` varchar(255) NOT NULL DEFAULT "",
  `port` int NOT NULL DEFAULT 0,
  `rpcPort` int NOT NULL DEFAULT 0,
  `createdAt` datetime NOT NULL,
  `updatedAt` datetime NOT NULL,
  PRIMARY KEY (`id`)
);

