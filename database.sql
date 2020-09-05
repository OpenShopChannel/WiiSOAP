/*

< Database Template File >
This file automatically adds the default database and tables.
WiiSOAP uses MySQL.

This SQL File does not guarantee functionality as WiiSOAP is still in early development statements.
It is suggested that you should hold off from using WiiSOAP unless you are confident that you know what you are doing.
Follow and practice proper security practices before handling user data.

*/

-- Generation Time: Jan 23, 2019 at 12:40 PM
-- Server version: 10.5.5-MariaDB
-- PHP Version: 7.3.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `owned_titles`
--

DROP TABLE IF EXISTS `owned_titles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `owned_titles` (
                                `account_id` varchar(9) NOT NULL,
                                `ticket_id` int(16) NOT NULL,
                                `title_id` varchar(16) NOT NULL,
                                `revocation_date` int(11) NOT NULL DEFAULT 0,
                                KEY `match_shop_title_metadata` (`title_id`),
                                KEY `order_account_ids` (`account_id`),
                                CONSTRAINT `match_shop_title_metadata` FOREIGN KEY (`title_id`) REFERENCES `shop_titles` (`title_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `shop_titles`
--

DROP TABLE IF EXISTS `shop_titles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `shop_titles` (
                               `title_id` varchar(16) NOT NULL,
                               `version` int(11) NOT NULL,
                               `description` mediumtext DEFAULT 'yada yada',
                               PRIMARY KEY (`title_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `userbase`
--

DROP TABLE IF EXISTS `userbase`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `userbase` (
                            `DeviceId` varchar(10) NOT NULL,
                            `DeviceToken` varchar(64) NOT NULL COMMENT 'This token should be considered a secret, so after generation only the sha256sum of the md5 the Wii sends is inserted.',
                            `AccountId` varchar(9) NOT NULL,
                            `Region` varchar(3) NOT NULL,
                            `Country` varchar(2) NOT NULL,
                            `Language` varchar(2) NOT NULL,
                            `SerialNo` varchar(11) NOT NULL,
                            `DeviceCode` varchar(16) NOT NULL,
                            PRIMARY KEY (`AccountId`),
                            UNIQUE KEY `AccountId` (`AccountId`),
                            UNIQUE KEY `userbase_DeviceId_uindex` (`DeviceId`),
                            UNIQUE KEY `userbase_DeviceToken_uindex` (`DeviceToken`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
