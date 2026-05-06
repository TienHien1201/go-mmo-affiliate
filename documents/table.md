-- Create database (run as a superuser or a role with CREATEDB)
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'mmmo_affiliate_db') THEN
		CREATE DATABASE mmmo_affiliate_db;
	END IF;
END $$;

-- Connect to the database (psql only)
-- \c mmmo_affiliate_db

-- Core tables
CREATE TABLE IF NOT EXISTS tenants (
id uuid PRIMARY KEY,
name varchar(255),
slug varchar(255),
plan varchar(255),
settings jsonb,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS users (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
username varchar(255),
email varchar(255),
full_name varchar(255),
avatar varchar(255),
password varchar(255),
password_hash varchar(255),
facebook_id varchar(255),
google_id varchar(255),
role_id bigint DEFAULT 2,
totp_secret varchar(255),
deleted_by bigint NOT NULL DEFAULT 0,
is_deleted boolean NOT NULL DEFAULT false,
deleted_at timestamptz,
"role" varchar(255),
is_active boolean,
last_login_at timestamptz,
created_at timestamptz,
updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roles (
id bigserial PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
name varchar(255),
description varchar(255),
is_active boolean DEFAULT true,
deleted_by bigint NOT NULL DEFAULT 0,
is_deleted boolean NOT NULL DEFAULT false,
deleted_at timestamptz,
created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS permissions (
id bigserial PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
name varchar(255) NOT NULL,
endpoint varchar(255) NOT NULL,
method varchar(100) NOT NULL,
module varchar(100) NOT NULL,
deleted_by bigint NOT NULL DEFAULT 0,
is_deleted boolean NOT NULL DEFAULT false,
deleted_at timestamptz,
created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
created_by bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS role_permission (
id bigserial PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
role_id bigint NOT NULL REFERENCES roles(id),
permission_id bigint NOT NULL REFERENCES permissions(id),
is_active boolean DEFAULT true,
deleted_by bigint NOT NULL DEFAULT 0,
is_deleted boolean NOT NULL DEFAULT false,
deleted_at timestamptz,
created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_role_permission_role_id ON role_permission(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permission_permission_id ON role_permission(permission_id);
CREATE INDEX IF NOT EXISTS idx_role_permission_tenant_id ON role_permission(tenant_id);
CREATE INDEX IF NOT EXISTS idx_role_permission_tenant_role ON role_permission(tenant_id, role_id);
CREATE INDEX IF NOT EXISTS idx_role_permission_tenant_permission ON role_permission(tenant_id, permission_id);

CREATE TABLE IF NOT EXISTS user_roles (
user_id uuid REFERENCES users(id),
role_id bigint REFERENCES roles(id),
assigned_at timestamptz,
PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
user_id uuid REFERENCES users(id),
action varchar(255),
entity_type varchar(255),
entity_id uuid,
payload jsonb,
ip_address varchar(255),
created_at timestamptz
);

-- Proxy context
CREATE TABLE IF NOT EXISTS proxy_pools (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
name varchar(255),
provider varchar(255),
proxy_type varchar(255),
status varchar(255),
total_count integer,
alive_count integer,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS proxies (
id uuid PRIMARY KEY,
pool_id uuid REFERENCES proxy_pools(id),
host varchar(255),
port integer,
username varchar(255),
password varchar(255),
protocol varchar(255),
status varchar(255),
country_code varchar(10),
city varchar(255),
latency_ms integer,
uptime_pct double precision,
fail_count integer,
last_checked_at timestamptz,
cooldown_until timestamptz,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS proxy_health_checks (
id uuid PRIMARY KEY,
proxy_id uuid REFERENCES proxies(id),
is_alive boolean,
latency_ms integer,
ip_resolved varchar(255),
error_message varchar(255),
checked_at timestamptz
);

-- Account context
CREATE TABLE IF NOT EXISTS device_fingerprints (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
canvas_hash varchar(255),
webgl_hash varchar(255),
screen_resolution varchar(255),
timezone varchar(255),
language varchar(255),
navigator_props jsonb,
is_used boolean,
used_by_account_id uuid,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS tiktok_accounts (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
assigned_user_id uuid REFERENCES users(id),
proxy_id uuid REFERENCES proxies(id),
device_fingerprint_id uuid REFERENCES device_fingerprints(id),
username varchar(255),
display_name varchar(255),
status varchar(255),
trust_score integer,
follower_count integer,
following_count integer,
total_gmv bigint,
shop_id varchar(255),
seller_level varchar(255),
cookies jsonb,
last_active_at timestamptz,
banned_at timestamptz,
created_at timestamptz,
updated_at timestamptz
);

ALTER TABLE device_fingerprints
ADD CONSTRAINT device_fingerprints_used_by_account_id_fkey
FOREIGN KEY (used_by_account_id) REFERENCES tiktok_accounts(id);

CREATE TABLE IF NOT EXISTS account_profiles (
id uuid PRIMARY KEY,
account_id uuid REFERENCES tiktok_accounts(id),
bio varchar(255),
avatar_url varchar(255),
region varchar(255),
language varchar(255),
categories jsonb,
sync_status varchar(255),
synced_at timestamptz
);

CREATE TABLE IF NOT EXISTS account_sessions (
id uuid PRIMARY KEY,
account_id uuid REFERENCES tiktok_accounts(id),
session_token text,
status varchar(255),
user_agent varchar(255),
platform varchar(255),
expires_at timestamptz,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS warmup_plans (
id uuid PRIMARY KEY,
account_id uuid REFERENCES tiktok_accounts(id),
status varchar(255),
duration_days integer,
current_day integer,
daily_targets jsonb,
started_at timestamptz,
completed_at timestamptz
);

CREATE TABLE IF NOT EXISTS warmup_logs (
id uuid PRIMARY KEY,
plan_id uuid REFERENCES warmup_plans(id),
account_id uuid REFERENCES tiktok_accounts(id),
day_number integer,
action_type varchar(255),
success boolean,
note varchar(255),
executed_at timestamptz
);

CREATE TABLE IF NOT EXISTS account_proxy_bindings (
id uuid PRIMARY KEY,
account_id uuid REFERENCES tiktok_accounts(id),
proxy_id uuid REFERENCES proxies(id),
is_active boolean,
bound_at timestamptz,
unbound_at timestamptz
);

-- Product context
CREATE TABLE IF NOT EXISTS product_templates (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
created_by uuid REFERENCES users(id),
name varchar(255),
category varchar(255),
base_title varchar(255),
base_description text,
base_hashtags jsonb,
base_price numeric(18, 2),
cost_price numeric(18, 2),
stock_quantity integer,
status varchar(255),
created_at timestamptz,
updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS product_listings (
id uuid PRIMARY KEY,
template_id uuid REFERENCES product_templates(id),
account_id uuid REFERENCES tiktok_accounts(id),
tiktok_product_id varchar(255),
title varchar(255),
description text,
hashtags jsonb,
price numeric(18, 2),
compare_price numeric(18, 2),
stock_quantity integer,
status varchar(255),
seo_score double precision,
view_count integer,
order_count integer,
revenue bigint,
published_at timestamptz,
created_at timestamptz,
updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS product_variants (
id uuid PRIMARY KEY,
listing_id uuid REFERENCES product_listings(id),
name varchar(255),
sku varchar(255),
price numeric(18, 2),
stock integer,
attributes jsonb,
status varchar(255),
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS product_media (
id uuid PRIMARY KEY,
listing_id uuid REFERENCES product_listings(id),
media_type varchar(255),
original_url varchar(255),
processed_url varchar(255),
gcs_path varchar(255),
sort_order integer,
status varchar(255),
metadata jsonb,
uploaded_at timestamptz
);

CREATE TABLE IF NOT EXISTS ai_generated_contents (
id uuid PRIMARY KEY,
listing_id uuid REFERENCES product_listings(id),
content_type varchar(255),
content_value text,
model_used varchar(255),
prompt_version varchar(255),
quality_score double precision,
is_active boolean,
generated_at timestamptz
);

CREATE TABLE IF NOT EXISTS seo_keywords (
id uuid PRIMARY KEY,
listing_id uuid REFERENCES product_listings(id),
keyword varchar(255),
search_volume integer,
competition_score double precision,
rank_position integer,
tracked_at timestamptz
);

CREATE TABLE IF NOT EXISTS price_history (
id uuid PRIMARY KEY,
listing_id uuid REFERENCES product_listings(id),
old_price numeric(18, 2),
new_price numeric(18, 2),
change_reason varchar(255),
changed_by uuid REFERENCES users(id),
changed_at timestamptz
);

-- Order context
CREATE TABLE IF NOT EXISTS orders (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
account_id uuid REFERENCES tiktok_accounts(id),
tiktok_order_id varchar(255),
status varchar(255),
subtotal numeric(18, 2),
shipping_fee numeric(18, 2),
discount numeric(18, 2),
total_amount numeric(18, 2),
buyer_name varchar(255),
buyer_phone varchar(255),
shipping_province varchar(255),
shipping_district varchar(255),
shipping_address varchar(255),
tracking_number varchar(255),
carrier_code varchar(255),
ordered_at timestamptz,
confirmed_at timestamptz,
shipped_at timestamptz,
delivered_at timestamptz,
created_at timestamptz,
updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS order_lines (
id uuid PRIMARY KEY,
order_id uuid REFERENCES orders(id),
listing_id uuid REFERENCES product_listings(id),
variant_id uuid REFERENCES product_variants(id),
product_name varchar(255),
variant_name varchar(255),
quantity integer,
unit_price numeric(18, 2),
total_price numeric(18, 2),
status varchar(255)
);

CREATE TABLE IF NOT EXISTS refund_requests (
id uuid PRIMARY KEY,
order_id uuid REFERENCES orders(id),
tiktok_refund_id varchar(255),
reason varchar(255),
status varchar(255),
refund_amount numeric(18, 2),
handler_note varchar(255),
handled_by uuid REFERENCES users(id),
requested_at timestamptz,
resolved_at timestamptz
);

CREATE TABLE IF NOT EXISTS fulfillment_logs (
id uuid PRIMARY KEY,
order_id uuid REFERENCES orders(id),
event_type varchar(255),
description varchar(255),
metadata jsonb,
actor_id uuid REFERENCES users(id),
logged_at timestamptz
);

-- Automation context
CREATE TABLE IF NOT EXISTS bot_scripts (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
name varchar(255),
action_type varchar(255),
script_content text,
version varchar(255),
status varchar(255),
language varchar(255),
created_by uuid REFERENCES users(id),
created_at timestamptz,
deprecated_at timestamptz
);

CREATE TABLE IF NOT EXISTS automation_campaigns (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
created_by uuid REFERENCES users(id),
bot_script_id uuid REFERENCES bot_scripts(id),
name varchar(255),
action_type varchar(255),
status varchar(255),
config jsonb,
schedule jsonb,
max_tasks_per_day integer,
completed_tasks integer,
failed_tasks integer,
started_at timestamptz,
stopped_at timestamptz,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS campaign_account_targets (
campaign_id uuid REFERENCES automation_campaigns(id),
account_id uuid REFERENCES tiktok_accounts(id),
status varchar(255),
tasks_done integer,
assigned_at timestamptz,
PRIMARY KEY (campaign_id, account_id)
);

CREATE TABLE IF NOT EXISTS automation_tasks (
id uuid PRIMARY KEY,
campaign_id uuid REFERENCES automation_campaigns(id),
account_id uuid REFERENCES tiktok_accounts(id),
action_type varchar(255),
status varchar(255),
input_params jsonb,
result jsonb,
retry_count integer,
error_message varchar(255),
scheduled_at timestamptz,
started_at timestamptz,
completed_at timestamptz
);

CREATE TABLE IF NOT EXISTS human_behavior_profiles (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
name varchar(255),
min_delay_ms integer,
max_delay_ms integer,
scroll_pattern jsonb,
click_pattern jsonb,
mistake_rate double precision,
is_active boolean,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS captcha_solver_configs (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
provider varchar(255),
api_key_encrypted varchar(255),
is_active boolean,
success_rate double precision,
avg_solve_ms integer,
created_at timestamptz
);

-- Analytics context
CREATE TABLE IF NOT EXISTS account_metrics (
id uuid PRIMARY KEY,
account_id uuid REFERENCES tiktok_accounts(id),
metric_date date,
views integer,
likes integer,
comments integer,
shares integer,
followers_gained integer,
followers_lost integer,
shop_visits integer,
gmv bigint,
orders_count integer,
conversion_rate double precision,
shop_score double precision,
recorded_at timestamptz
);

CREATE TABLE IF NOT EXISTS product_metrics (
id uuid PRIMARY KEY,
listing_id uuid REFERENCES product_listings(id),
metric_date date,
impressions integer,
clicks integer,
ctr double precision,
orders_count integer,
revenue bigint,
refunds_count integer,
avg_rating double precision,
reviews_count integer,
recorded_at timestamptz
);

CREATE TABLE IF NOT EXISTS revenue_snapshots (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
period_type varchar(255),
period_start date,
period_end date,
total_gmv bigint,
total_orders integer,
avg_order_value numeric(18, 2),
total_refunds bigint,
net_revenue bigint,
estimated_profit bigint,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS trend_signals (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
keyword varchar(255),
category varchar(255),
trend_score double precision,
search_volume integer,
growth_rate double precision,
region varchar(255),
related_products jsonb,
status varchar(255),
detected_at timestamptz,
expires_at timestamptz
);

CREATE TABLE IF NOT EXISTS performance_reports (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
report_type varchar(255),
status varchar(255),
period_start date,
period_end date,
summary_data jsonb,
gcs_export_path varchar(255),
generated_by uuid REFERENCES users(id),
generated_at timestamptz
);

-- Finance context
CREATE TABLE IF NOT EXISTS cost_entries (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
account_id uuid REFERENCES tiktok_accounts(id),
campaign_id uuid REFERENCES automation_campaigns(id),
cost_type varchar(255),
description varchar(255),
amount numeric(18, 2),
currency varchar(10),
cost_date date,
recorded_by uuid REFERENCES users(id),
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS commission_records (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
user_id uuid REFERENCES users(id),
period varchar(255),
total_gmv bigint,
commission_rate numeric(10, 4),
commission_amount numeric(18, 2),
status varchar(255),
approved_by uuid REFERENCES users(id),
approved_at timestamptz,
created_at timestamptz
);

-- Notification context
CREATE TABLE IF NOT EXISTS notification_rules (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
event_type varchar(255),
channel varchar(255),
severity varchar(255),
conditions jsonb,
template jsonb,
is_active boolean,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS notifications (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
rule_id uuid REFERENCES notification_rules(id),
user_id uuid REFERENCES users(id),
channel varchar(255),
title varchar(255),
body text,
status varchar(255),
metadata jsonb,
sent_at timestamptz,
read_at timestamptz
);

CREATE TABLE IF NOT EXISTS webhooks (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
url varchar(255),
events jsonb,
secret_hash varchar(255),
is_active boolean,
fail_count integer,
last_triggered_at timestamptz,
created_at timestamptz
);

CREATE TABLE IF NOT EXISTS domain_events (
id uuid PRIMARY KEY,
tenant_id uuid REFERENCES tenants(id),
event_type varchar(255),
aggregate_type varchar(255),
aggregate_id uuid,
payload jsonb,
is_processed boolean,
retry_count integer,
occurred_at timestamptz,
processed_at timestamptz
);

-- Compliance context
CREATE TABLE IF NOT EXISTS account_risk_scores (
id uuid PRIMARY KEY,
account_id uuid REFERENCES tiktok_accounts(id),
risk_score double precision,
risk_level varchar(255),
risk_factors jsonb,
recommended_action varchar(255),
evaluated_at timestamptz
);

CREATE TABLE IF NOT EXISTS violation_logs (
id uuid PRIMARY KEY,
account_id uuid REFERENCES tiktok_accounts(id),
violation_type varchar(255),
severity varchar(255),
description varchar(255),
action_taken varchar(255),
detected_at timestamptz,
resolved_at timestamptz
);

-- ===============================
-- FOREIGN KEYS
-- ===============================
ALTER TABLE users
ADD CONSTRAINT fk_users_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE roles
ADD CONSTRAINT fk_roles_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE permissions
ADD CONSTRAINT fk_permissions_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE role_permission
ADD CONSTRAINT fk_role_permission_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE user_roles
ADD CONSTRAINT fk_user_roles_user
FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE user_roles
ADD CONSTRAINT fk_user_roles_role
FOREIGN KEY (role_id) REFERENCES roles(id);

ALTER TABLE audit_logs
ADD CONSTRAINT fk_audit_logs_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE audit_logs
ADD CONSTRAINT fk_audit_logs_user
FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE tiktok_accounts
ADD CONSTRAINT fk_tiktok_accounts_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE tiktok_accounts
ADD CONSTRAINT fk_tiktok_accounts_user
FOREIGN KEY (assigned_user_id) REFERENCES users(id);

ALTER TABLE tiktok_accounts
ADD CONSTRAINT fk_tiktok_accounts_proxy
FOREIGN KEY (proxy_id) REFERENCES proxies(id);

ALTER TABLE tiktok_accounts
ADD CONSTRAINT fk_tiktok_accounts_device
FOREIGN KEY (device_fingerprint_id) REFERENCES device_fingerprints(id);

ALTER TABLE account_profiles
ADD CONSTRAINT fk_account_profiles_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE account_sessions
ADD CONSTRAINT fk_account_sessions_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE device_fingerprints
ADD CONSTRAINT fk_device_fingerprints_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE device_fingerprints
ADD CONSTRAINT fk_device_fingerprints_used_by
FOREIGN KEY (used_by_account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE warmup_plans
ADD CONSTRAINT fk_warmup_plans_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE warmup_logs
ADD CONSTRAINT fk_warmup_logs_plan
FOREIGN KEY (plan_id) REFERENCES warmup_plans(id);

ALTER TABLE warmup_logs
ADD CONSTRAINT fk_warmup_logs_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE proxy_pools
ADD CONSTRAINT fk_proxy_pools_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE proxies
ADD CONSTRAINT fk_proxies_pool
FOREIGN KEY (pool_id) REFERENCES proxy_pools(id);

ALTER TABLE proxy_health_checks
ADD CONSTRAINT fk_proxy_health_checks_proxy
FOREIGN KEY (proxy_id) REFERENCES proxies(id);

ALTER TABLE account_proxy_bindings
ADD CONSTRAINT fk_account_proxy_bindings_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE account_proxy_bindings
ADD CONSTRAINT fk_account_proxy_bindings_proxy
FOREIGN KEY (proxy_id) REFERENCES proxies(id);

ALTER TABLE product_templates
ADD CONSTRAINT fk_product_templates_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE product_templates
ADD CONSTRAINT fk_product_templates_created_by
FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE product_listings
ADD CONSTRAINT fk_product_listings_template
FOREIGN KEY (template_id) REFERENCES product_templates(id);

ALTER TABLE product_listings
ADD CONSTRAINT fk_product_listings_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE product_variants
ADD CONSTRAINT fk_product_variants_listing
FOREIGN KEY (listing_id) REFERENCES product_listings(id);

ALTER TABLE product_media
ADD CONSTRAINT fk_product_media_listing
FOREIGN KEY (listing_id) REFERENCES product_listings(id);

ALTER TABLE ai_generated_contents
ADD CONSTRAINT fk_ai_generated_contents_listing
FOREIGN KEY (listing_id) REFERENCES product_listings(id);

ALTER TABLE seo_keywords
ADD CONSTRAINT fk_seo_keywords_listing
FOREIGN KEY (listing_id) REFERENCES product_listings(id);

ALTER TABLE price_history
ADD CONSTRAINT fk_price_history_listing
FOREIGN KEY (listing_id) REFERENCES product_listings(id);

ALTER TABLE price_history
ADD CONSTRAINT fk_price_history_changed_by
FOREIGN KEY (changed_by) REFERENCES users(id);

ALTER TABLE orders
ADD CONSTRAINT fk_orders_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE orders
ADD CONSTRAINT fk_orders_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE order_lines
ADD CONSTRAINT fk_order_lines_order
FOREIGN KEY (order_id) REFERENCES orders(id);

ALTER TABLE order_lines
ADD CONSTRAINT fk_order_lines_listing
FOREIGN KEY (listing_id) REFERENCES product_listings(id);

ALTER TABLE order_lines
ADD CONSTRAINT fk_order_lines_variant
FOREIGN KEY (variant_id) REFERENCES product_variants(id);

ALTER TABLE refund_requests
ADD CONSTRAINT fk_refund_requests_order
FOREIGN KEY (order_id) REFERENCES orders(id);

ALTER TABLE refund_requests
ADD CONSTRAINT fk_refund_requests_handler
FOREIGN KEY (handled_by) REFERENCES users(id);

ALTER TABLE fulfillment_logs
ADD CONSTRAINT fk_fulfillment_logs_order
FOREIGN KEY (order_id) REFERENCES orders(id);

ALTER TABLE fulfillment_logs
ADD CONSTRAINT fk_fulfillment_logs_actor
FOREIGN KEY (actor_id) REFERENCES users(id);

ALTER TABLE bot_scripts
ADD CONSTRAINT fk_bot_scripts_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE bot_scripts
ADD CONSTRAINT fk_bot_scripts_created_by
FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE automation_campaigns
ADD CONSTRAINT fk_automation_campaigns_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE automation_campaigns
ADD CONSTRAINT fk_automation_campaigns_created_by
FOREIGN KEY (created_by) REFERENCES users(id);

ALTER TABLE automation_campaigns
ADD CONSTRAINT fk_automation_campaigns_bot_script
FOREIGN KEY (bot_script_id) REFERENCES bot_scripts(id);

ALTER TABLE campaign_account_targets
ADD CONSTRAINT fk_campaign_account_targets_campaign
FOREIGN KEY (campaign_id) REFERENCES automation_campaigns(id);

ALTER TABLE campaign_account_targets
ADD CONSTRAINT fk_campaign_account_targets_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE automation_tasks
ADD CONSTRAINT fk_automation_tasks_campaign
FOREIGN KEY (campaign_id) REFERENCES automation_campaigns(id);

ALTER TABLE automation_tasks
ADD CONSTRAINT fk_automation_tasks_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE human_behavior_profiles
ADD CONSTRAINT fk_human_behavior_profiles_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE captcha_solver_configs
ADD CONSTRAINT fk_captcha_solver_configs_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE account_metrics
ADD CONSTRAINT fk_account_metrics_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE product_metrics
ADD CONSTRAINT fk_product_metrics_listing
FOREIGN KEY (listing_id) REFERENCES product_listings(id);

ALTER TABLE revenue_snapshots
ADD CONSTRAINT fk_revenue_snapshots_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE trend_signals
ADD CONSTRAINT fk_trend_signals_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE performance_reports
ADD CONSTRAINT fk_performance_reports_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE performance_reports
ADD CONSTRAINT fk_performance_reports_generated_by
FOREIGN KEY (generated_by) REFERENCES users(id);

ALTER TABLE cost_entries
ADD CONSTRAINT fk_cost_entries_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE cost_entries
ADD CONSTRAINT fk_cost_entries_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE cost_entries
ADD CONSTRAINT fk_cost_entries_campaign
FOREIGN KEY (campaign_id) REFERENCES automation_campaigns(id);

ALTER TABLE cost_entries
ADD CONSTRAINT fk_cost_entries_recorded_by
FOREIGN KEY (recorded_by) REFERENCES users(id);

ALTER TABLE commission_records
ADD CONSTRAINT fk_commission_records_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE commission_records
ADD CONSTRAINT fk_commission_records_user
FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE commission_records
ADD CONSTRAINT fk_commission_records_approved_by
FOREIGN KEY (approved_by) REFERENCES users(id);

ALTER TABLE notification_rules
ADD CONSTRAINT fk_notification_rules_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE notifications
ADD CONSTRAINT fk_notifications_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE notifications
ADD CONSTRAINT fk_notifications_rule
FOREIGN KEY (rule_id) REFERENCES notification_rules(id);

ALTER TABLE notifications
ADD CONSTRAINT fk_notifications_user
FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE webhooks
ADD CONSTRAINT fk_webhooks_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE domain_events
ADD CONSTRAINT fk_domain_events_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE account_risk_scores
ADD CONSTRAINT fk_account_risk_scores_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);

ALTER TABLE violation_logs
ADD CONSTRAINT fk_violation_logs_account
FOREIGN KEY (account_id) REFERENCES tiktok_accounts(id);
