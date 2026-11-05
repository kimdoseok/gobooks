# PL/PgSQL
## Basic function
```
CREATE OR REPLACE FUNCTION add_numbers(a INT, b INT)
RETURNS INT AS $$
DECLARE
    sum_result INT;
BEGIN
    sum_result := a + b;
    RETURN sum_result;
END;
$$ LANGUAGE plpgsql;

-- To run it:
SELECT add_numbers(10, 25);

```

## Controls
```
CREATE OR REPLACE FUNCTION check_shipping_tier(order_amount NUMERIC)
RETURNS TEXT AS $$
BEGIN
    IF order_amount >= 100.00 THEN
        RETURN 'Free Express Shipping';
    ELSIF order_amount >= 50.00 THEN
        RETURN 'Standard Shipping';
    ELSE
        RETURN 'Paid Shipping Required';
    END IF;
END;
$$ LANGUAGE plpgsql;

```

## Looping
```
CREATE OR REPLACE FUNCTION process_inactive_users()
RETURNS VOID AS $$
DECLARE
    user_record RECORD; -- Generic variable to hold a single row
BEGIN
    FOR user_record IN 
        SELECT id, email FROM users WHERE last_login < NOW() - INTERVAL '1 year'
    LOOP
        -- Perform logic for each user row
        RAISE NOTICE 'Processing inactive user: %', user_record.email;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

```

## Building blocks
```
* Variables & Assignments: Variables must be declared before use and assigned values using the := operator.
* Control Structures: Supports standard conditionals (IF, THEN, ELSE, CASE) and loop structures (LOOP, WHILE, FOR).
* Dynamic SQL: Allows execution of strings as SQL commands using the EXECUTE statement when table or column names need to change dynamically.
* Error Handling: Uses BEGIN ... EXCEPTION WHEN ... THEN blocks to trap runtime exceptions gracefully and roll back problematic statements.
```