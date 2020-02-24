/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 80019
 Source Host           : 127.0.0.1:3306
 Source Schema         : avn_wallet

 Target Server Type    : MySQL
 Target Server Version : 80019
 File Encoding         : 65001

 Date: 24/02/2020 12:21:41
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for transactions
-- ----------------------------
DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions` (
  `id` binary(16) NOT NULL,
  `walletID` binary(16) NOT NULL,
  `balance` decimal(10,2) NOT NULL,
  `createdAt` datetime NOT NULL,
  `cause` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '',
  `description` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  `type` enum('PROMOTION','MANUALLY_INCREASE','MANUALLY_DECREASE','TRANSFER','SYSTEM_DECREASE') COLLATE utf8mb4_general_ci NOT NULL,
  `status` enum('PENDING','ACTIVE') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'PENDING',
  PRIMARY KEY (`id`),
  KEY `walletID` (`walletID`),
  CONSTRAINT `walletID` FOREIGN KEY (`walletID`) REFERENCES `wallet` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` binary(16) NOT NULL,
  `firstName` varchar(70) COLLATE utf8mb4_general_ci DEFAULT '',
  `lastName` varchar(70) COLLATE utf8mb4_general_ci DEFAULT '',
  `cellphone` bigint unsigned NOT NULL,
  `email` varchar(100) COLLATE utf8mb4_general_ci DEFAULT '',
  `status` enum('ACTIVE','PENDING','INACTIVE','DELETED') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'ACTIVE',
  `createdAt` datetime NOT NULL,
  `updatedAt` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `cellphone` (`cellphone`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for wallet
-- ----------------------------
DROP TABLE IF EXISTS `wallet`;
CREATE TABLE `wallet` (
  `id` binary(16) NOT NULL,
  `userID` binary(16) NOT NULL,
  `charge` decimal(10,2) unsigned NOT NULL,
  `createdAt` datetime NOT NULL,
  `updateAt` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `userID` (`userID`),
  CONSTRAINT `userID` FOREIGN KEY (`userID`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Procedure structure for getTransactions
-- ----------------------------
DROP PROCEDURE IF EXISTS `getTransactions`;
delimiter ;;
CREATE PROCEDURE `avn_wallet`.`getTransactions`(IN CELLPHONE BIGINT(10))
BEGIN
  
	SELECT
	BIN_TO_UUID(t.`id`) AS id,
	t.`balance` AS balance,
	t.`type` AS type,
	t.`createdAt` AS createdAt,
	t.`cause` AS cause,
	t.`description` AS description
	FROM `transactions` AS t
	INNER JOIN `wallet` AS w ON w.`id` = t.`walletID`
	INNER JOIN `users` AS u ON w.`userID` = u.`id`
	WHERE u.`cellphone` = CELLPHONE
	AND t.`status` = 'ACTIVE'
	ORDER BY t.`createdAt` DESC;

END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for getWallet
-- ----------------------------
DROP PROCEDURE IF EXISTS `getWallet`;
delimiter ;;
CREATE PROCEDURE `avn_wallet`.`getWallet`(IN CELLPHONE BIGINT(10))
BEGIN

	SELECT
		w.charge 
	FROM
		users AS us
		INNER JOIN wallet AS w ON us.id = w.userID 
	WHERE
		us.cellphone = CELLPHONE;
	
END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for insertWallet
-- ----------------------------
DROP PROCEDURE IF EXISTS `insertWallet`;
delimiter ;;
CREATE PROCEDURE `avn_wallet`.`insertWallet`(IN CHARGE_DATA BIGINT(15), IN CELLPHONE_DATA BIGINT(10))
BEGIN
  	
	DECLARE WALLET_ID BINARY(16) DEFAULT NULL;
	
	DECLARE USER_ID BINARY(16) DEFAULT NULL;
	
  DECLARE EXIT HANDLER FOR SQLEXCEPTION
		BEGIN
		ROLLBACK;
		RESIGNAL;
		END;
	
	DECLARE EXIT HANDLER FOR SQLWARNING
		BEGIN
		ROLLBACK;
		RESIGNAL;
		END;
	
		##select user
		SET USER_ID := (SELECT u.`id` FROM `users` as u WHERE u.`cellphone` = CELLPHONE_DATA limit 1);

		##select wallet
		SET WALLET_ID := (SELECT w.`id` FROM `wallet` as w WHERE w.`userID` = USER_ID);

		
		START TRANSACTION;
					
          IF(WALLET_ID IS NULL) THEN
						    
								SET WALLET_ID := UUID_TO_BIN(UUID());
								
								##insert into wallet
								INSERT INTO `wallet` (`id`, `charge`, `userID`, `createdAt`)
								VALUES (WALLET_ID, 0, USER_ID, NOW());
					
					END IF;
										
	
					#TODO the data can come from producer call
					##insert transaction
					INSERT INTO `transactions` (`id`,`walletID`,`balance`,`type`,`createdAt`)
					VALUES (UUID_TO_BIN(UUID()),WALLET_ID,CHARGE_DATA,'PROMOTION',NOW());
					
					
		COMMIT;
	
		SELECT w.`charge`,us.`firstName`,us.`lastName`,us.`cellphone` FROM `wallet` AS w INNER JOIN `users` AS us ON w.`userID`=us.`id` WHERE w.`userID` = USER_ID;

END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for seeder
-- ----------------------------
DROP PROCEDURE IF EXISTS `seeder`;
delimiter ;;
CREATE PROCEDURE `avn_wallet`.`seeder`()
BEGIN
 
	INSERT INTO `users` (`id`, `firstName`,`lastName`,`cellphone`,`email`,`status`,`createdAt`)
	VALUES
	(uuid_to_bin(uuid()),"amirali","roshanaei",9118000217,"mr.roshanae@gmail.com","ACTIVE",NOW()),
	(uuid_to_bin(uuid()),"ali","hosseini",9123457689,"ali.hosseini@yahoo.com","ACTIVE",NOW()),
	(uuid_to_bin(uuid()),"mohsen","majidi",9378675432,"mohsen.majidi@gmail.com","ACTIVE",NOW()),
	(uuid_to_bin(uuid()),"reza","mahdavi",9219453298,"mr.roshanae@hotmail.com","ACTIVE",NOW()),
	(uuid_to_bin(uuid()),"farhad","jamshidi",9159674312,"farhad.jamshidi@outlook.com","ACTIVE",NOW());

END
;;
delimiter ;

-- ----------------------------
-- Procedure structure for updateWallet
-- ----------------------------
DROP PROCEDURE IF EXISTS `updateWallet`;
delimiter ;;
CREATE PROCEDURE `avn_wallet`.`updateWallet`(IN CELLPHONE_DATA BIGINT(10))
BEGIN
  	
	DECLARE WALLET_ID BINARY(16) DEFAULT NULL;
	
	DECLARE USER_ID BINARY(16) DEFAULT NULL;
	
	DECLARE CHARGE_DATA DECIMAL(10,2) DEFAULT NULL;
	
  DECLARE EXIT HANDLER FOR SQLEXCEPTION
		BEGIN
		ROLLBACK;
		RESIGNAL;
		END;
	
	DECLARE EXIT HANDLER FOR SQLWARNING
		BEGIN
		ROLLBACK;
		RESIGNAL;
		END;
	
		##select user
		SET USER_ID := (SELECT u.`id` FROM `users` AS u WHERE u.`cellphone` = CELLPHONE_DATA LIMIT 1);

		##select wallet
		SET WALLET_ID := (SELECT w.`id` FROM `wallet` AS w WHERE w.`userID` = USER_ID LIMIT 1);

		##select charge
		SET CHARGE_DATA := (SELECT t.`balance` FROM `transactions` AS t WHERE t.`walletID` = WALLET_ID ORDER BY t.`createdAt` DESC LIMIT 1);
		
		START TRANSACTION;
				
					##update transaction
					UPDATE `transactions` AS t
					SET t.`status` = 'ACTIVE'
					WHERE t.`walletID` = WALLET_ID
					ORDER BY t.`createdAt` DESC 
					LIMIT 1;
					
					##update wallet
					UPDATE `wallet` AS w
					SET w.`charge` = w.`charge` + CHARGE_DATA
					WHERE w.`id` = WALLET_ID;
					
					
		COMMIT;
	
		SELECT w.`charge`,us.`firstName`,us.`lastName`,us.`cellphone` FROM `wallet` AS w INNER JOIN `users` AS us ON w.`userID`=us.`id` WHERE w.`userID` = USER_ID;

END
;;
delimiter ;

SET FOREIGN_KEY_CHECKS = 1;
