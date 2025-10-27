-- reverse: set comment to table: "posts"
COMMENT ON TABLE "posts" IS '';
-- reverse: create index "post_title_user_posts" to table: "posts"
DROP INDEX "post_title_user_posts";
-- reverse: create "posts" table
DROP TABLE "posts";
-- reverse: set comment to table: "users"
COMMENT ON TABLE "users" IS '';
-- reverse: create index "users_name_key" to table: "users"
DROP INDEX "users_name_key";
-- reverse: create index "users_email_key" to table: "users"
DROP INDEX "users_email_key";
-- reverse: create index "user_name_email" to table: "users"
DROP INDEX "user_name_email";
-- reverse: create "users" table
DROP TABLE "users";
