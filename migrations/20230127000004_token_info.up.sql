CREATE TABLE IF NOT EXISTS "token_info" (
  "id" varchar(1024) NOT NULL,
  "name" varchar(512) NOT NULL,
  "website" varchar(1024) NULL,
  "description" text,
  "explorer" varchar(1024) NULL,
  "type" varchar(512) NULL,
  "blockchain" varchar(255) NULL,
  "symbol" varchar(255) NULL,
  "status" varchar(255) NULL,  
  "decimals" int NOT NULL DEFAULT '0'
) 
 