-- MySQL dump 10.13  Distrib 5.7.31, for macos10.14 (x86_64)
--
-- Host: localhost    Database: cloudops
-- ------------------------------------------------------
-- Server version	5.7.31

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `apis`
--

DROP TABLE IF EXISTS `apis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `apis` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `path` varchar(50) DEFAULT NULL COMMENT '路由路径',
  `method` varchar(50) DEFAULT NULL COMMENT 'http请求方法',
  `pid` bigint(20) DEFAULT NULL COMMENT 'apiGroups 父级的id 为了给树用的',
  `title` varchar(50) DEFAULT NULL COMMENT '名称',
  `type` varchar(5) DEFAULT NULL COMMENT '类型 0=父级 1=子级',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_apis_title` (`title`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `apis`
--

LOCK TABLES `apis` WRITE;
/*!40000 ALTER TABLE `apis` DISABLE KEYS */;
INSERT INTO `apis` VALUES (1,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/*','ALL',0,'所有-写操作','0'),(2,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/*','GET',0,'所有-读操作','0'),(3,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/system/*','ALL',0,'系统管理-写操作','0'),(4,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/system/*','GET',0,'系统管理-读操作','0'),(5,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/stree/*','ALL',0,'服务树-写操作','0'),(6,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/stree/*','GET',0,'服务树-读操作','0'),(7,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/workorder/*','ALL',0,'工单管理-写操作','0'),(8,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/workorder/*','GET',0,'工单管理-读操作','0'),(9,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/jobexec/*','ALL',0,'任务执行中心-写操作','0'),(10,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/jobexec/*','GET',0,'任务执行中心-读操作','0'),(11,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/monitor/*','ALL',0,'Prometheus监控-写操作','0'),(12,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/monitor/*','GET',0,'Prometheus监控-读操作','0'),(13,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/k8s/*','ALL',0,'k8s管理员-写操作','0'),(14,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/k8s/*','GET',0,'k8s多集群-读操作','0'),(15,'2024-07-24 21:50:48.496','2024-07-24 21:50:48.496','/api/k8sApp/*','ALL',0,'k8s应用-写操作','0');
/*!40000 ALTER TABLE `apis` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `casbin_rule`
--

DROP TABLE IF EXISTS `casbin_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `casbin_rule` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `casbin_rule`
--

LOCK TABLES `casbin_rule` WRITE;
/*!40000 ALTER TABLE `casbin_rule` DISABLE KEYS */;
INSERT INTO `casbin_rule` VALUES (28,'p','dev','/api/*','DELETE','','',''),(15,'p','dev','/api/*','GET','','',''),(24,'p','dev','/api/*','POST','','',''),(29,'p','dev','/api/jobexec/*','DELETE','','',''),(26,'p','dev','/api/jobexec/*','POST','','',''),(27,'p','dev','/api/monitor/*','DELETE','','',''),(23,'p','dev','/api/monitor/*','POST','','',''),(22,'p','dev','/api/stree/*','DELETE','','',''),(20,'p','dev','/api/stree/*','POST','','',''),(17,'p','dev','/api/system/*','DELETE','','',''),(16,'p','dev','/api/system/*','POST','','',''),(25,'p','dev','/api/workorder/*','DELETE','','',''),(21,'p','dev','/api/workorder/*','POST','','',''),(9,'p','k8s_admin','/api/*','GET','','',''),(11,'p','k8s_admin','/api/k8s/*','DELETE','','',''),(8,'p','k8s_admin','/api/k8s/*','GET','','',''),(10,'p','k8s_admin','/api/k8s/*','POST','','',''),(7,'p','ops','/api/*','GET','','',''),(6,'p','ops','/api/k8sApp/*','DELETE','','',''),(4,'p','ops','/api/k8sApp/*','GET','','',''),(5,'p','ops','/api/k8sApp/*','POST','','',''),(3,'p','super','/api/*','DELETE','','',''),(1,'p','super','/api/*','GET','','',''),(2,'p','super','/api/*','POST','','',''),(12,'p','test','/api/*','GET','','',''),(14,'p','test','/api/stree/*','DELETE','','',''),(13,'p','test','/api/stree/*','POST','','',''),(19,'p','test','/api/system/*','DELETE','','',''),(18,'p','test','/api/system/*','POST','','','');
/*!40000 ALTER TABLE `casbin_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `menus`
--

DROP TABLE IF EXISTS `menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `menus` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL COMMENT '名称',
  `title` longtext COMMENT '名称',
  `pid` bigint(20) DEFAULT NULL COMMENT '父级的id',
  `parent_menu` longtext COMMENT '父级的id',
  `icon` longtext COMMENT '图标',
  `type` varchar(5) DEFAULT NULL COMMENT '类型 0=目录 1=子菜单',
  `show` varchar(5) DEFAULT NULL COMMENT '类型 0=禁用 1=启用',
  `order_no` bigint(20) DEFAULT NULL COMMENT '排序',
  `component` varchar(50) DEFAULT NULL COMMENT '前端组件 菜单就是LAYOUT',
  `redirect` varchar(50) DEFAULT NULL COMMENT '显示路径',
  `path` varchar(50) DEFAULT NULL COMMENT '路由路径',
  `remark` longtext COMMENT '用户描述',
  `home_path` longtext COMMENT '登陆后的默认首页',
  `status` varchar(191) DEFAULT '1' COMMENT '是否启用 0禁用 1启用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_menus_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=68 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `menus`
--

LOCK TABLES `menus` WRITE;
/*!40000 ALTER TABLE `menus` DISABLE KEYS */;
INSERT INTO `menus` VALUES (1,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','System','系统管理',0,'','ion:settings-outline','0','1',1,'LAYOUT','/system/account','/system','','','1'),(2,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','Permission','权限管理',0,'','ion:layers-outline','0','0',2,'LAYOUT','/permission/front/page','/permission','','','1'),(3,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','ServiceTree','服务树与cmdb',0,'','ion:layers-outline','0','1',3,'LAYOUT','/serviceTree/serviceTree/index','/serviceTree','','','1'),(4,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','WorkOrder','工单系统',0,'','ion:git-compare-outline','0','1',4,'LAYOUT','/workOrder/process/index','/workOrder','','','1'),(5,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','JobExec','任务执行中心',0,'','ion:git-compare-outline','0','1',5,'LAYOUT','/jobExec/task/index','/jobExec','','','1'),(6,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','PrometheusMonitor','prometheus监控',0,'','ion:tv-outline','0','1',6,'LAYOUT','/monitor/scrape/index','/monitor','','','1'),(7,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sAdmins','[管理员]k8s多集群管理',0,'','ant-design:cloud-server-outlined','0','1',7,'LAYOUT','/k8s/cluster/index','/k8s','','','1'),(8,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sOpsManagement','[运维]k8s应用管理',0,'','ant-design:dropbox-square-filled','0','1',8,'LAYOUT','/k8sApplication/project/index','/k8sApplication','','','1'),(9,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','CicdManagement','cicd管理',0,'','ant-design:printer-filled','0','1',9,'LAYOUT','/cicd/workorder/index','/cicd','','','1'),(10,'2024-07-24 21:50:48.489','2024-08-10 14:26:41.136','MenuManagement','菜单管理',1,'','ant-design:appstore-twotone','1','1',1,'/demo/system/menu/index','','menu','','','1'),(11,'2024-07-24 21:50:48.489','2024-08-10 14:27:19.083','AccountManagement','用户管理',1,'','ant-design:user-outlined','1','1',2,'/demo/system/account/index','','account','','','1'),(12,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','RoleManagement','角色管理',1,'','ion:layers-outline','1','1',3,'/demo/system/role/index','','role','','','1'),(13,'2024-07-24 21:50:48.489','2024-08-09 22:47:37.247','ChangePassword','修改密码',1,'','ion:layers-outline','1','1',4,'/demo/system/password/index','','changePassword','','','1'),(14,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','ApiManagement','api接口管理',1,'','ant-design:account-book-filled','1','1',5,'/demo/system/api/index','','api','','','1'),(15,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','PermissionFrontDemo','前端权限管理',2,'','ion:layers-outline','1','1',1,'/demo/permission/front/index','','front','','','1'),(16,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','ServiceTreeIndex','服务树',3,'','ant-design:account-book-filled','1','1',1,'/stree/stree/index','','stree','','','1'),(17,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','ServiceTreeIndexAsync','服务树异步',3,'','ant-design:account-book-filled','1','1',2,'/stree/stree/indexAsync','','streeAsync','','','1'),(18,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','ProcessManagement','流程管理',4,'','ant-design:account-book-filled','1','1',1,'/workorder/process/index','','process','','','1'),(19,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','FormDesignManagement','表单设计管理',4,'','ant-design:account-book-filled','1','1',3,'/workorder/formDesign/index','','formDesign','','','1'),(20,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','WorkOrderTemplateManagement','工单模板管理',4,'','ant-design:account-book-filled','1','1',4,'/workorder/template/index','','template','','','1'),(21,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','WorkOrderApply','工单申请',4,'','ant-design:account-book-filled','1','1',4,'/workorder/apply/index','','apply','','','1'),(22,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','WorkOrderCreate','工单创建',4,'','ant-design:account-book-filled','1','0',4,'/workorder/apply/create','','create','','','1'),(23,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','WorkOrderSearch','工单查询',4,'','simple-icons:about-dot-me','1','1',4,'/workorder/apply/search','','search','','','1'),(24,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','WorkOrderDetail','工单详情',4,'','ant-design:account-book-filled','1','0',4,'/workorder/detail/index','','detail','','','1'),(25,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','JobExecScript','脚本管理',5,'','ant-design:account-book-filled','1','1',4,'/jobexec/script/index','','script','','','1'),(26,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','JobExecTask','任务管理',5,'','ant-design:code-filled','1','1',5,'/jobexec/task/index','','task','','','1'),(27,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','JobExecDetail','任务详情',5,'','ant-design:account-book-filled','1','0',6,'/jobexec/detail/index','','detail','','','1'),(28,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorPool','Prometheus集群管理',6,'','ant-design:undo-outlined','1','1',1,'/monitor/pool/index','','pool','','','1'),(29,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorScrapeJob','prometheus采集job',6,'','ant-design:file-search-outlined','1','1',2,'/monitor/scrape/index','','scrape','','','1'),(30,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorAlertmanagerPool','alertmanager集群管理',6,'','ant-design:alert-outlined','1','1',3,'/monitor/alertmanager/index','','alertmanager','','','1'),(31,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorOndutyGroup','值班组管理',6,'','ant-design:calendar-outlined','1','1',4,'/monitor/ondutygroup/index','','ondutygroup','','','1'),(32,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorSendGroup','发送组管理',6,'','ant-design:dingding-outlined','1','1',5,'/monitor/sendgroup/index','','sendgroup','','','1'),(33,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorAlertRule','告警规则管理',6,'','ant-design:sisternode-outlined','1','1',6,'/monitor/alertrule/index','','alertrule','','','1'),(34,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorAlertEvent','告警事件管理',6,'','ant-design:project-outlined','1','1',7,'/monitor/event/index','','event','','','1'),(35,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorRecordRule','预聚合规则管理',6,'','ant-design:heat-map-outlined','1','1',8,'/monitor/recordrule/index','','recordrule','','','1'),(36,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorOndutyGroupPlan','值班组排班',6,'','ant-design:calendar-outlined','1','0',20,'/monitor/ondutygroup/plan','','ondutygroup-plan','','','1'),(37,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorPrometheusMainYaml','prometheus主配置查看',6,'','ant-design:undo-outlined','1','0',20,'/monitor/pool/mainYaml','','mainYaml','','','1'),(38,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorPrometheusAlertRuleYaml','prometheus告警规则配置查看',6,'','ant-design:undo-outlined','1','0',20,'/monitor/pool/alertRuleYaml','','alertRuleYaml','','','1'),(39,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorPrometheusRecordRuleYaml','prometheus预聚合规则配置查看',6,'','ant-design:undo-outlined','1','0',20,'/monitor/pool/recordRuleYaml','','recordRuleYaml','','','1'),(40,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','MonitorAlertmanagerYaml','Alertmanager配置查看',6,'','ant-design:undo-outlined','1','0',20,'/monitor/alertmanager/yaml','','alertManagerYaml','','','1'),(41,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sClustersManagement','k8s集群管理',7,'','ion:layers-outline','1','1',1,'/k8s/cluster/index','','cluster','','','1'),(42,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sNodesManagement','k8s节点管理',7,'','ant-design:hdd-outlined','1','1',2,'/k8s/node/index','','node','','','1'),(43,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sNodesDetail','k8s节点管理',7,'','ant-design:insert-row-above-outlined','1','0',2,'/k8s/node/detail','','nodeDetail','','','1'),(44,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sYamlTemplate','k8s-yaml模板管理',7,'','ant-design:expand-outlined','1','1',3,'/k8s/yaml/template','','yamlTemplate','','','1'),(45,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sYamlTask','k8s-yaml任务管理',7,'','ant-design:merge-cells-outlined','1','1',4,'/k8s/yaml/task','','yamlTask','','','1'),(46,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sPodManagement','k8s-pod管理',7,'','ant-design:cluster-outlined','1','1',5,'/k8s/podAdmin/index','','podList','','','1'),(47,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sPodYaml','k8s-pod-yaml展示',7,'','ant-design:shop-two-tone','1','0',5,'/k8s/podAdmin/podYaml','','podYaml','','','1'),(48,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sDeployYaml','k8s-Deploy-yaml展示',7,'','ant-design:shop-two-tone','1','0',5,'/k8s/deployAdmin/deployYaml','','deployYaml','','','1'),(49,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sConfigmapYaml','k8s-configmap-yaml展示',7,'','ant-design:shop-two-tone','1','0',5,'/k8s/configmapAdmin/configmapYaml','','configmapYaml','','','1'),(50,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sServiceYaml','k8s-service-yaml展示',7,'','ant-design:shop-two-tone','1','0',5,'/k8s/serviceAdmin/serviceYaml','','serviceYaml','','','1'),(51,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sDeployManagement','k8s-deployment管理',7,'','ant-design:mac-command-filled','1','1',6,'/k8s/deployAdmin/index','','deployList','','','1'),(52,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sConfigmapManagement','k8s-Configmap管理',7,'','ant-design:folder-open-filled','1','1',7,'/k8s/configmapAdmin/index','','configmapList','','','1'),(53,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sSvcManagement','k8s-svc管理',7,'','ant-design:api-filled','1','1',7,'/k8s/serviceAdmin/index','','serviceList','','','1'),(54,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sPodNormalLog','k8s-pod非tail-log',7,'','ant-design:shop-two-tone','1','0',5,'/k8s/podAdmin/normalLog','','podNormalLog','','','1'),(55,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sPodTailLog','k8s-pod-ws-tail-log',7,'','ant-design:shop-two-tone','1','0',5,'/k8s/podAdmin/tailLog','','podTailLog','','','1'),(56,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sPodWebShell','k8s-pod命令执行',7,'','ant-design:shop-two-tone','1','0',5,'/k8s/podAdmin/webshell','','podWebShell','','','1'),(57,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sProjectsManagement','k8s项目管理',8,'','ant-design:project-filled','1','1',1,'/k8sApplication/project/index','','project','','','1'),(58,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sApplicationsManagement','k8s应用管理',8,'','ant-design:apple-filled','1','1',2,'/k8sApplication/application/index','','application','','','1'),(59,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sApplicationsEditor','k8s应用编辑',8,'','ant-design:apple-filled','1','0',2,'/k8sApplication/application/editor','','editor','','','1'),(60,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sInstancesManagement','k8s实例管理',8,'','ant-design:codepen-square-filled','1','1',3,'/k8sApplication/instance/index','','instance','','','1'),(61,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sCronjobManagement','k8s计划任务管理',8,'','ant-design:schedule-filled','1','1',4,'/k8sApplication/cronjob/index','','cronjob','','','1'),(62,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sInstancesEditor','k8s实例编辑',8,'','ant-design:apple-filled','1','0',2,'/k8sApplication/instance/editor','','instanceEditor','','','1'),(63,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sInstancesLogExec','k8s实例日志和exec',8,'','ant-design:apple-filled','1','0',2,'/k8sApplication/instance/logExec','','instanceLogExec','','','1'),(64,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','K8sCronjobsEditor','k8s计划任务编辑',8,'','ant-design:apple-filled','1','0',2,'/k8sApplication/cronjob/editor','','cronjobEditor','','','1'),(65,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','CicdWorkOrderManagement','cicd工单列表',9,'','ant-design:funnel-plot-filled','1','1',1,'/cicd/workorder/index','','orderList','','','1'),(66,'2024-07-24 21:50:48.489','2024-07-24 21:50:48.489','CicdWorkOrderDetail','cicd发布单详情',9,'','ant-design:funnel-plot-filled','1','0',1,'/cicd/workorder/detail','','orderDetail','','','1');
/*!40000 ALTER TABLE `menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ops_admins`
--

DROP TABLE IF EXISTS `ops_admins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ops_admins` (
  `stree_node_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`stree_node_id`,`user_id`),
  KEY `fk_ops_admins_user` (`user_id`),
  CONSTRAINT `fk_ops_admins_stree_node` FOREIGN KEY (`stree_node_id`) REFERENCES `stree_nodes` (`id`),
  CONSTRAINT `fk_ops_admins_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ops_admins`
--

LOCK TABLES `ops_admins` WRITE;
/*!40000 ALTER TABLE `ops_admins` DISABLE KEYS */;
INSERT INTO `ops_admins` VALUES (1,1),(2,1),(1,2),(2,2);
/*!40000 ALTER TABLE `ops_admins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rd_admins`
--

DROP TABLE IF EXISTS `rd_admins`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rd_admins` (
  `stree_node_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`stree_node_id`,`user_id`),
  KEY `fk_rd_admins_user` (`user_id`),
  CONSTRAINT `fk_rd_admins_stree_node` FOREIGN KEY (`stree_node_id`) REFERENCES `stree_nodes` (`id`),
  CONSTRAINT `fk_rd_admins_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rd_admins`
--

LOCK TABLES `rd_admins` WRITE;
/*!40000 ALTER TABLE `rd_admins` DISABLE KEYS */;
INSERT INTO `rd_admins` VALUES (6,1),(8,1),(6,2),(8,2);
/*!40000 ALTER TABLE `rd_admins` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `rd_members`
--

DROP TABLE IF EXISTS `rd_members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rd_members` (
  `stree_node_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`stree_node_id`,`user_id`),
  KEY `fk_rd_members_user` (`user_id`),
  CONSTRAINT `fk_rd_members_stree_node` FOREIGN KEY (`stree_node_id`) REFERENCES `stree_nodes` (`id`),
  CONSTRAINT `fk_rd_members_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rd_members`
--

LOCK TABLES `rd_members` WRITE;
/*!40000 ALTER TABLE `rd_members` DISABLE KEYS */;
/*!40000 ALTER TABLE `rd_members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_apis`
--

DROP TABLE IF EXISTS `role_apis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role_apis` (
  `api_id` bigint(20) unsigned NOT NULL,
  `role_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`api_id`,`role_id`),
  KEY `fk_role_apis_role` (`role_id`),
  CONSTRAINT `fk_role_apis_api` FOREIGN KEY (`api_id`) REFERENCES `apis` (`id`),
  CONSTRAINT `fk_role_apis_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_apis`
--

LOCK TABLES `role_apis` WRITE;
/*!40000 ALTER TABLE `role_apis` DISABLE KEYS */;
INSERT INTO `role_apis` VALUES (1,1),(2,1),(3,1),(4,1),(5,1),(6,1),(7,1),(8,1),(9,1),(10,1),(11,1),(12,1),(13,1),(14,1),(15,1),(2,2),(15,2),(2,4),(13,4),(2,8),(3,8),(4,8),(2,9);
/*!40000 ALTER TABLE `role_apis` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_menus`
--

DROP TABLE IF EXISTS `role_menus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role_menus` (
  `menu_id` bigint(20) unsigned NOT NULL,
  `role_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`menu_id`,`role_id`),
  KEY `fk_role_menus_role` (`role_id`),
  CONSTRAINT `fk_role_menus_menu` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`),
  CONSTRAINT `fk_role_menus_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_menus`
--

LOCK TABLES `role_menus` WRITE;
/*!40000 ALTER TABLE `role_menus` DISABLE KEYS */;
INSERT INTO `role_menus` VALUES (1,1),(2,1),(3,1),(4,1),(5,1),(6,1),(7,1),(8,1),(9,1),(10,1),(11,1),(12,1),(13,1),(14,1),(15,1),(16,1),(17,1),(18,1),(19,1),(20,1),(21,1),(22,1),(23,1),(24,1),(25,1),(26,1),(27,1),(28,1),(29,1),(30,1),(31,1),(32,1),(33,1),(34,1),(35,1),(36,1),(37,1),(38,1),(39,1),(40,1),(41,1),(42,1),(43,1),(44,1),(45,1),(46,1),(47,1),(48,1),(49,1),(50,1),(51,1),(52,1),(53,1),(54,1),(55,1),(56,1),(57,1),(58,1),(59,1),(60,1),(61,1),(62,1),(63,1),(64,1),(65,1),(66,1),(13,2),(16,4),(17,4),(18,4),(19,4),(20,4),(21,4),(22,4),(23,4),(24,4),(41,4),(42,4),(43,4),(44,4),(45,4),(46,4),(47,4),(48,4),(49,4),(50,4),(51,4),(52,4),(53,4),(10,8),(11,8),(12,8),(13,8),(14,8),(13,9);
/*!40000 ALTER TABLE `role_menus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `roles` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `order_no` bigint(20) DEFAULT NULL COMMENT '排序',
  `role_name` varchar(100) DEFAULT NULL COMMENT '角色中文名称',
  `role_value` varchar(100) DEFAULT NULL COMMENT '角色值',
  `remark` longtext COMMENT '用户描述',
  `home_path` longtext COMMENT '登陆后的默认首页',
  `status` varchar(191) DEFAULT '1' COMMENT '角色是否被冻结 1正常 2冻结',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_roles_role_name` (`role_name`),
  UNIQUE KEY `idx_roles_role_value` (`role_value`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'2024-07-24 21:50:48.487','2024-07-24 21:50:48.487',0,'超级管理员','super','','','1'),(2,'2024-07-24 21:50:48.501','2024-08-09 22:56:26.172',0,'普通运维','ops','','','1'),(3,'2024-07-24 21:50:48.507','2024-07-24 21:50:48.507',0,'后台机器人','bot_super','','','1'),(4,'2024-07-24 21:50:48.509','2024-08-03 16:09:06.266',0,'k8s集群管理员','k8s_admin','','','1'),(8,'2024-08-03 16:25:20.509','2024-08-09 22:55:41.644',0,'测试用户','test',NULL,'','1'),(9,'2024-08-03 16:25:44.820','2024-08-10 09:21:24.404',0,'开发','dev','','','1');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `stree_nodes`
--

DROP TABLE IF EXISTS `stree_nodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `stree_nodes` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `title` varchar(50) DEFAULT NULL COMMENT '名称',
  `pid` bigint(20) unsigned DEFAULT NULL COMMENT 'StreeNodeGroups 父级的id 为了给树用的',
  `level` bigint(20) DEFAULT NULL COMMENT '层级 ',
  `is_leaf` tinyint(1) DEFAULT NULL COMMENT '是否启用 0=否 1=是',
  `desc` longtext COMMENT '描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `pid_title` (`title`,`pid`),
  KEY `idx_stree_nodes_pid` (`pid`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `stree_nodes`
--

LOCK TABLES `stree_nodes` WRITE;
/*!40000 ALTER TABLE `stree_nodes` DISABLE KEYS */;
INSERT INTO `stree_nodes` VALUES (16,'2024-08-11 10:09:44.722','2024-08-11 10:09:44.722','腾讯',0,1,0,''),(21,'2024-08-11 16:30:22.165','2024-08-11 16:30:22.165','微信',16,2,0,''),(22,'2024-08-11 16:31:40.802','2024-08-11 16:31:40.802','QQ',0,1,0,''),(23,'2024-08-11 16:32:06.842','2024-08-11 16:32:06.842','元梦之星',22,2,0,'');
/*!40000 ALTER TABLE `stree_nodes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_roles`
--

DROP TABLE IF EXISTS `user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_roles` (
  `role_id` bigint(20) unsigned NOT NULL,
  `user_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`role_id`,`user_id`),
  KEY `fk_user_roles_user` (`user_id`),
  CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_roles`
--

LOCK TABLES `user_roles` WRITE;
/*!40000 ALTER TABLE `user_roles` DISABLE KEYS */;
INSERT INTO `user_roles` VALUES (1,1),(2,1),(8,2),(3,3),(4,4),(9,16);
/*!40000 ALTER TABLE `user_roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL COMMENT '用户id',
  `username` varchar(100) DEFAULT NULL COMMENT '用户登录名',
  `password` longtext COMMENT '用户登录密码',
  `real_name` varchar(100) DEFAULT NULL COMMENT '用户昵称 中文名',
  `desc` longtext COMMENT '用户描述',
  `mobile` longtext COMMENT '手机号',
  `fei_shu_user_id` longtext COMMENT '飞书userid  获取方式请看 https://open.feishu.cn/document/home/user-identity-introduction/open-id',
  `account_type` bigint(20) DEFAULT '1' COMMENT '用户是否是服务账号 1普通用户 2服务账号',
  `service_account_token` longtext,
  `home_path` longtext COMMENT '登陆后的默认首页',
  `enable` bigint(20) DEFAULT '1' COMMENT '用户是否被冻结 1正常 2冻结',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_real_name` (`real_name`),
  UNIQUE KEY `idx_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'2024-07-24 21:50:48.486','2024-07-27 20:26:31.909',0,'admin','$2a$10$BWCka/k2FSLHXurmEngJ6Oe2jxhywuzle7UOn.4XDidSbHSVVzH2q','花海','','15810947075','6d31gdf4',1,'','/system/account',1),(2,'2024-07-24 21:50:48.500','2024-08-09 22:57:17.983',0,'test','$2a$10$HmvS0hXZICjg2eajOKhB5O1PHsWJiqRJbaQKAmiaulFozCfXb5YOq','狗子','','15810947075','38g18781',1,'','/system/account',1),(3,'2024-07-24 21:50:48.507','2024-07-24 21:50:48.507',0,'auto_order_robot','$2a$10$sC9FwqZf/t3EQcA6NxH/GuvX7DIrHyqMeZUdbIwi6U/GTwerjejjC','自动工单执行机器人','','15810947075','6d31gdf4',1,'','/system/account',1),(4,'2024-07-24 21:50:48.509','2024-07-24 21:50:48.509',0,'k8s_admin_01','$2a$10$bntI4mzlizk/7o7O68245OJwfFl9rWZHlszd5LiFdtJ2zAm5p8gru','k8s管理员','','15810947075','6d31gdf4',1,'','/k8s/cluster',1),(16,'2024-08-07 10:44:29.807','2024-08-07 10:44:29.807',0,'cicd','123456','二狗子','','','666666',1,'','',1);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-08-11 16:34:28
