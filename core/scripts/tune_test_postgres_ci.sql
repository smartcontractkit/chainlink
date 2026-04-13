-- Applied after chainlink_test exists (see ci-core.yml). Database-level defaults for test throughput.
ALTER DATABASE chainlink_test SET synchronous_commit TO off;
ALTER DATABASE chainlink_test SET jit TO off;
ALTER DATABASE chainlink_test SET random_page_cost TO 1.1;
ALTER DATABASE chainlink_test SET effective_cache_size TO '1GB';
