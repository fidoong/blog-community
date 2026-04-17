-- 批量生成真实风格测试数据。
-- 默认会在现有数据之后继续插入，避免破坏既有结构与内容。
-- 生成内容：用户、文章、评论、关注、点赞、收藏，并回填计数。

BEGIN;

CREATE TEMP TABLE seed_meta (
    prefix text not null,
    password_hash text not null
) ON COMMIT DROP;

INSERT INTO seed_meta(prefix, password_hash)
VALUES ('seed20260417', '$2a$10$iqWGee1hqgRRkK1E6nZR2eqFdzRPp9YeSBCIN6WrkS932CXrOWvBi');

INSERT INTO users (email, username, password_hash, avatar_url, oauth_provider, oauth_id, role, created_at, updated_at)
SELECT
    lower(m.prefix) || '_' || gs::text || '@example.com' AS email,
    lower(m.prefix) || '_user_' || lpad(gs::text, 4, '0') AS username,
    m.password_hash,
    'https://api.dicebear.com/7.x/notionists/svg?seed=' || gs::text AS avatar_url,
    CASE WHEN gs % 5 = 0 THEN 'google' WHEN gs % 7 = 0 THEN 'github' ELSE 'none' END AS oauth_provider,
    CASE WHEN gs % 5 = 0 OR gs % 7 = 0 THEN lower(m.prefix) || '_oauth_' || gs::text ELSE NULL END AS oauth_id,
    CASE WHEN gs % 20 = 0 THEN 'admin' ELSE 'user' END AS role,
    now() - ((180 - gs % 180) || ' days')::interval - ((gs % 23) || ' hours')::interval AS created_at,
    now() - ((gs % 17) || ' hours')::interval AS updated_at
FROM seed_meta m
CROSS JOIN generate_series(1, 240) AS gs
ON CONFLICT (email) DO NOTHING;

CREATE TEMP TABLE seed_user_pool ON COMMIT DROP AS
SELECT id, username
FROM users
WHERE username LIKE 'seed20260417_user_%';

CREATE TEMP TABLE seed_title_pool (
    idx int not null,
    title text not null
) ON COMMIT DROP;

INSERT INTO seed_title_pool(idx, title)
VALUES
    (1, 'Go 服务上线前的十个自检动作'),
    (2, '把博客系统拆成模块化单体之后，我学到什么'),
    (3, 'Redis 缓存穿透排查手记'),
    (4, '一次真实的 PostgreSQL 索引优化记录'),
    (5, '从零搭一个内容社区，需要哪些后端基础能力'),
    (6, '评论系统如何设计楼中楼结构'),
    (7, 'JWT 与 Refresh Token 的边界在哪里'),
    (8, '前后端分离项目的本地开发协作规范'),
    (9, '我如何给个人博客增加推荐流'),
    (10, '把接口错误码体系做统一之后的收益'),
    (11, 'Gin 中间件链排错技巧'),
    (12, '真实项目里怎么设计点赞收藏表'),
    (13, '内容审核状态流转的工程化实现'),
    (14, '为什么我选择 Ent 作为 Go ORM'),
    (15, '把需求拆成领域模块后的代码结构变化'),
    (16, '一个周末做完博客后台管理台的经验'),
    (17, '如何设计更像真人写的技术文章摘要'),
    (18, 'ElasticSearch 接入社区搜索的最小闭环'),
    (19, '一次前端登录注册联调中的踩坑复盘'),
    (20, 'Docker Compose 管理本地开发基础设施');

INSERT INTO posts (title, content, summary, content_type, cover_image, author_id, status, view_count, like_count, comment_count, collect_count, tags, published_at, created_at, updated_at)
SELECT
    tp.title || CASE WHEN gs > 1 THEN '（案例 ' || gs::text || '）' ELSE '' END AS title,
    '## ' || tp.title || E'\n\n'
        || '这是一篇围绕真实项目经验整理的文章，记录问题背景、定位过程与最终方案。' || E'\n\n'
        || '### 背景' || E'\n'
        || '最近在维护内容社区项目时，我重新梳理了服务拆分、缓存设计、接口契约和数据一致性问题。' || E'\n\n'
        || '### 过程' || E'\n'
        || '我先从日志和监控入手，再逐步下钻到数据库、缓存和业务层。过程中发现很多问题并非代码错误，而是配置、约束和边界条件没有被完整表达。' || E'\n\n'
        || '### 结论' || E'\n'
        || '对工程系统来说，稳定的默认值、明确的错误语义，以及可重复执行的验证脚本，比一次性的修补更重要。' AS content,
    '围绕 ' || tp.title || ' 的实战总结，覆盖设计取舍、排障过程与最终落地方案。' AS summary,
    CASE WHEN (tp.idx + gs) % 6 = 0 THEN 'rich_text' ELSE 'markdown' END AS content_type,
    'https://images.unsplash.com/photo-' || (1500000000000 + tp.idx * 1000 + gs)::text || '?auto=format&fit=crop&w=1200&q=80' AS cover_image,
    au.id AS author_id,
    CASE
        WHEN (tp.idx + gs) % 11 = 0 THEN 'draft'
        WHEN (tp.idx + gs) % 13 = 0 THEN 'pending'
        WHEN (tp.idx + gs) % 17 = 0 THEN 'rejected'
        ELSE 'published'
    END AS status,
    0,
    0,
    0,
    0,
    to_jsonb(ARRAY[
        (ARRAY['Go','PostgreSQL','Redis','Gin','Ent','架构','性能优化','工程实践','前端协作','接口设计'])[((tp.idx + 1) % 10) + 1],
        (ARRAY['社区产品','微服务','缓存','数据库','日志追踪','部署','认证鉴权','评论系统','推荐流','测试'])[((tp.idx + 3) % 10) + 1],
        (ARRAY['排障','设计模式','数据建模','CI/CD','开发体验','中间件','后端','全栈','内容平台','可观测性'])[((tp.idx + 5) % 10) + 1]
    ]) AS tags,
    CASE
        WHEN ((tp.idx + gs) % 11 = 0 OR (tp.idx + gs) % 13 = 0 OR (tp.idx + gs) % 17 = 0) THEN NULL
        ELSE now() - ((tp.idx * 3 + gs - 2) || ' days')::interval
    END AS published_at,
    now() - ((tp.idx * 3 + gs) || ' days')::interval AS created_at,
    now() - ((tp.idx * 3 + gs - 1) || ' days')::interval AS updated_at
FROM seed_title_pool tp
CROSS JOIN generate_series(1, 9) AS gs
JOIN LATERAL (
    SELECT id
    FROM seed_user_pool
    ORDER BY id
    OFFSET ((tp.idx * 13 + gs * 7) % (SELECT count(*) FROM seed_user_pool))
    LIMIT 1
) au ON true;

CREATE TEMP TABLE seed_post_pool ON COMMIT DROP AS
SELECT id, author_id, status
FROM posts
WHERE title LIKE 'Go 服务上线前的十个自检动作%'
   OR title LIKE '把博客系统拆成模块化单体之后，我学到什么%'
   OR title LIKE 'Redis 缓存穿透排查手记%'
   OR title LIKE '一次真实的 PostgreSQL 索引优化记录%'
   OR title LIKE '从零搭一个内容社区，需要哪些后端基础能力%'
   OR title LIKE '评论系统如何设计楼中楼结构%'
   OR title LIKE 'JWT 与 Refresh Token 的边界在哪里%'
   OR title LIKE '前后端分离项目的本地开发协作规范%'
   OR title LIKE '我如何给个人博客增加推荐流%'
   OR title LIKE '把接口错误码体系做统一之后的收益%'
   OR title LIKE 'Gin 中间件链排错技巧%'
   OR title LIKE '真实项目里怎么设计点赞收藏表%'
   OR title LIKE '内容审核状态流转的工程化实现%'
   OR title LIKE '为什么我选择 Ent 作为 Go ORM%'
   OR title LIKE '把需求拆成领域模块后的代码结构变化%'
   OR title LIKE '一个周末做完博客后台管理台的经验%'
   OR title LIKE '如何设计更像真人写的技术文章摘要%'
   OR title LIKE 'ElasticSearch 接入社区搜索的最小闭环%'
   OR title LIKE '一次前端登录注册联调中的踩坑复盘%'
   OR title LIKE 'Docker Compose 管理本地开发基础设施%';

CREATE TEMP TABLE seed_comment_roots ON COMMIT DROP AS
WITH inserted AS (
    INSERT INTO comments (content, post_id, author_id, parent_id, like_count, created_at, updated_at)
    SELECT
        CASE ((p.id + u.id + gs) % 8)
            WHEN 0 THEN '这篇写得很透，尤其是排查路径，照着做就能复现问题。'
            WHEN 1 THEN '我也遇到过类似情况，最后卡在配置和环境差异上。'
            WHEN 2 THEN '如果能把监控截图和关键 SQL 一起附上，信息会更完整。'
            WHEN 3 THEN '赞同，工程里的稳定性很多时候取决于默认值和边界条件。'
            WHEN 4 THEN '这段关于数据模型的解释很清楚，适合刚入门的同学。'
            WHEN 5 THEN '我们团队也在用类似方案，不过缓存失效策略不太一样。'
            WHEN 6 THEN '请问评论和点赞计数是同步更新还是异步聚合？'
            ELSE '内容很实在，没有空话，读完直接能拿去改自己的项目。'
        END AS content,
        p.id AS post_id,
        u.id AS author_id,
        NULL::bigint AS parent_id,
        0,
        now() - (((p.id + gs) % 45) || ' days')::interval - (((u.id + gs) % 18) || ' hours')::interval AS created_at,
        now() - (((p.id + gs) % 44) || ' days')::interval - (((u.id + gs) % 17) || ' hours')::interval AS updated_at
    FROM seed_post_pool p
    JOIN LATERAL (
        SELECT id
        FROM seed_user_pool
        WHERE id <> p.author_id
        ORDER BY id
        LIMIT 12
    ) u ON true
    JOIN generate_series(1, 2) AS gs ON true
    WHERE p.status = 'published'
    RETURNING id, post_id, author_id
)
SELECT id, post_id, author_id
FROM inserted;

INSERT INTO comments (content, post_id, author_id, parent_id, like_count, created_at, updated_at)
SELECT
    CASE (c.id % 6)
        WHEN 0 THEN '这个补充很关键，尤其是对新手来说。'
        WHEN 1 THEN '我也建议把这一步写进自动化脚本里。'
        WHEN 2 THEN '有类似实践，确实能明显减少联调成本。'
        WHEN 3 THEN '如果加上失败重试策略，整体会更稳。'
        WHEN 4 THEN '感谢分享，这个视角比单讲代码更有用。'
        ELSE '收到，我准备按你的思路重构一下现有逻辑。'
    END,
    c.post_id,
    u.id,
    c.id,
    0,
    now() - (((c.id + u.id) % 20) || ' days')::interval,
    now() - (((c.id + u.id) % 19) || ' days')::interval
FROM seed_comment_roots c
JOIN LATERAL (
    SELECT id
    FROM seed_user_pool
    WHERE id <> c.author_id
    ORDER BY id
    OFFSET (c.id % (SELECT count(*) FROM seed_user_pool))
    LIMIT 1
) u ON true
WHERE c.id % 3 = 0;

INSERT INTO follows (follower_id, following_id, created_at)
SELECT DISTINCT u1.id, u2.id, now() - (((u1.id + u2.id) % 120) || ' days')::interval
FROM seed_user_pool u1
JOIN seed_user_pool u2 ON u1.id <> u2.id
WHERE (u1.id + u2.id) % 19 = 0
ON CONFLICT (follower_id, following_id) DO NOTHING;

INSERT INTO like_records (target_type, target_id, user_id, created_at)
SELECT 'post', p.id, u.id, now() - (((p.id + u.id) % 60) || ' days')::interval
FROM seed_post_pool p
JOIN seed_user_pool u ON u.id <> p.author_id
WHERE p.status = 'published' AND (p.id * 7 + u.id * 11) % 23 = 0
ON CONFLICT (target_type, target_id, user_id) DO NOTHING;

INSERT INTO like_records (target_type, target_id, user_id, created_at)
SELECT 'comment', c.id, u.id, now() - (((c.id + u.id) % 40) || ' days')::interval
FROM comments c
JOIN seed_post_pool p ON p.id = c.post_id
JOIN seed_user_pool u ON u.id <> c.author_id
WHERE (c.id * 5 + u.id * 3) % 29 = 0
ON CONFLICT (target_type, target_id, user_id) DO NOTHING;

INSERT INTO collect_records (target_type, target_id, user_id, created_at)
SELECT 'post', p.id, u.id, now() - (((p.id + u.id) % 50) || ' days')::interval
FROM seed_post_pool p
JOIN seed_user_pool u ON u.id <> p.author_id
WHERE p.status = 'published' AND (p.id * 13 + u.id * 17) % 31 = 0
ON CONFLICT (target_type, target_id, user_id) DO NOTHING;

UPDATE comments c
SET like_count = COALESCE(src.like_count, 0),
    updated_at = GREATEST(c.updated_at, now() - ((c.id % 7) || ' hours')::interval)
FROM (
    SELECT target_id, count(*)::int AS like_count
    FROM like_records
    WHERE target_type = 'comment'
    GROUP BY target_id
) src
WHERE c.id = src.target_id;

UPDATE posts p
SET like_count = COALESCE(l.likes, 0),
    comment_count = COALESCE(cm.comments, 0),
    collect_count = COALESCE(ct.collects, 0),
    view_count = GREATEST(COALESCE(l.likes, 0) * 18 + COALESCE(cm.comments, 0) * 9 + (p.id % 500), 50),
    updated_at = GREATEST(p.updated_at, now() - ((p.id % 11) || ' hours')::interval)
FROM seed_post_pool seeded
LEFT JOIN (
    SELECT target_id, count(*)::int AS likes
    FROM like_records
    WHERE target_type = 'post'
    GROUP BY target_id
) l ON l.target_id = seeded.id
LEFT JOIN (
    SELECT post_id, count(*)::int AS comments
    FROM comments
    GROUP BY post_id
) cm ON cm.post_id = seeded.id
LEFT JOIN (
    SELECT target_id, count(*)::int AS collects
    FROM collect_records
    WHERE target_type = 'post'
    GROUP BY target_id
) ct ON ct.target_id = seeded.id
WHERE p.id = seeded.id;

COMMIT;
