DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id          INT AUTO_INCREMENT NOT NULL,
  uuid        BINARY(16) NOT NULL UNIQUE,
  name        VARCHAR(255),
  email       VARCHAR(255) NOT NULL UNIQUE,
  pwdmd5      CHAR(32) NOT NULL,
  salt        CHAR(64) NOT NULL,
  created_at  TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (`id`)
);
CREATE INDEX users_email ON users(email);

DROP TABLE IF EXISTS threads;
CREATE TABLE threads (
  id          INT AUTO_INCREMENT NOT NULL,
  uuid        BINARY(16) NOT NULL UNIQUE,
  topic_id    INT NOT NULL,
  title       VARCHAR(255) NOT NULL,
  content     VARCHAR(1000) NOT NULL,
  user_id     INT NOT NULL,
  last_pos    INT DEFAULT 1 NOT NULL,
  view_count  INT DEFAULT 0 NOT NULL,
  num_posts   INT DEFAULT 0 NOT NULL,
  created_at  TIMESTAMP DEFAULT NOW(),    
  PRIMARY KEY (`id`)
);
CREATE INDEX threads_userid ON threads(user_id);
CREATE INDEX threads_viewcount ON threads(view_count DESC);
CREATE INDEX threads_numposts ON threads(num_posts DESC);
CREATE INDEX threads_createdat ON threads(created_at DESC);

DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
  id          INT AUTO_INCREMENT NOT NULL,
  uuid        BINARY(16) NOT NULL UNIQUE,
  content     VARCHAR(500) NOT NULL,
  user_id     INT NOT NULL,  
  thread_id   INT NOT NULL,
  thread_pos  INT NOT NULL,
  created_at  TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (`id`)
);
CREATE INDEX posts_threadid_threadpos ON posts(thread_id, thread_pos ASC);
CREATE INDEX posts_userid ON posts(user_id);
CREATE INDEX posts_createdat ON posts(created_at DESC);

DROP TABLE IF EXISTS comments;
CREATE TABLE comments (
  id                INT AUTO_INCREMENT NOT NULL,
  uuid              BINARY(16) NOT NULL UNIQUE,  
  content           VARCHAR(300) NOT NULL,
  user_id           INT NOT NULL,
  reply_to_user_id  INT NOT NULL,
  reply_to_uuid     BINARY(16) NOT NULL,
  post_uuid         BINARY(16) NOT NULL,
  thread_uuid       BINARY(16) NOT NULL,
  created_at        TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (`id`)
);
CREATE INDEX comments_userid ON comments(user_id);
CREATE INDEX comments_postuuid ON comments(post_uuid);
CREATE INDEX comments_threaduuid ON comments(thread_uuid);
CREATE INDEX comments_createdat ON comments(created_at DESC);

DROP TABLE IF EXISTS topics;
CREATE TABLE topics (
  id          INT AUTO_INCREMENT NOT NULL,
  cate        VARCHAR(255) NOT NULL,    
  info        TEXT NOT NULL,  
  PRIMARY KEY (`id`)
);
INSERT INTO topics (cate, info)
VALUES ("News", "Current events, updates, and reports from around the world, delivered in a timely and objective manner."),
("Entertainment","Features on movies, music, TV shows, celebrities, and other forms of popular culture."),
("Education"," Content designed to inform, educate, and enhance knowledge in various subjects, including academics, skills, and life hacks."),
("Lifestyle"," Focuses on fashion, health, fitness, food, travel, and other aspects of daily life."),
("Technology","Covers the latest advancements, trends, and reviews in the world of technology, including gadgets, software, and the internet."),
("Finance","News, analysis, and insights on the global economy, markets, companies, and business strategies."),
("Sports","Covers various sports events, teams, players, and related news, opinions, and analysis."),
("Health and Wellness","Focuses on maintaining and improving physical and mental health, including fitness, nutrition, and mental health resources."),
("Arts and Culture","Features on fine arts, history, culture, and heritage, promoting understanding and appreciation of diverse artistic expressions and traditions."),
("Travel","Provides information on destinations, attractions, travel tips, and cultural experiences, aiming to inspire and assist in planning travel adventures.")

