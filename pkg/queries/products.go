package queries

const ListProducts = `
    SELECT 
        p.id,
        p.product_name,
        p.product_url,
        p.product_image_url,
        p.product_description,
        p.created_datetime,
        p.updated_datetime,
        uc.id,
        uc.email,
        uc.password_hash,
        uu.id,
        uu.email,
        uu.password_hash
    FROM products p
    JOIN users uc ON uc.id = p.created_by_user_id
    LEFT JOIN users uu ON uu.id = p.updated_by_user_id
    ORDER BY p.created_datetime DESC
`

const GetProductByID = `
    SELECT 
        p.id,
        p.product_name,
        p.product_url,
        p.product_image_url,
        p.product_description,
        p.created_datetime,
        p.updated_datetime,
        uc.id,
        uc.email,
        uc.password_hash,
        uu.id,
        uu.email,
        uu.password_hash
    FROM products p
    JOIN users uc ON uc.id = p.created_by_user_id
    LEFT JOIN users uu ON uu.id = p.updated_by_user_id
    WHERE p.id = $1
`

const InsertProductWithUsersCTE = `
    WITH ins AS (
        INSERT INTO products (product_name, product_url, product_image_url, product_description, created_by_user_id, updated_by_user_id)
        VALUES (
            $1,
            $2,
            $3,
            $4,
            NULLIF($5, '')::uuid,
            NULLIF($6, '')::uuid
        )
        RETURNING id, product_name, product_url, product_image_url, product_description, created_datetime, updated_datetime, created_by_user_id, updated_by_user_id
    )
    SELECT 
        ins.id,
        ins.product_name,
        ins.product_url,
        ins.product_image_url,
        ins.product_description,
        ins.created_datetime,
        ins.updated_datetime,
        uc.id,
        uc.email,
        uc.password_hash,
        uu.id,
        uu.email,
        uu.password_hash
    FROM ins
    JOIN users uc ON uc.id = ins.created_by_user_id
    LEFT JOIN users uu ON uu.id = ins.updated_by_user_id
`

const InsertProductWithoutUpdatedCTE = `
    WITH ins AS (
        INSERT INTO products (product_name, product_url, product_image_url, product_description, created_by_user_id)
        VALUES (
            $1,
            $2,
            $3,
            $4,
            NULLIF($5, '')::uuid
        )
        RETURNING id, product_name, product_url, product_image_url, product_description, created_datetime, updated_datetime, created_by_user_id, updated_by_user_id
    )
    SELECT 
        ins.id,
        ins.product_name,
        ins.product_url,
        ins.product_image_url,
        ins.product_description,
        ins.created_datetime,
        ins.updated_datetime,
        uc.id,
        uc.email,
        uc.password_hash,
        uu.id,
        uu.email,
        uu.password_hash
    FROM ins
    JOIN users uc ON uc.id = ins.created_by_user_id
    LEFT JOIN users uu ON uu.id = ins.updated_by_user_id
`

const UpdateProductWithUsersCTE = `
    WITH upd AS (
        UPDATE products
        SET product_name = $1,
            product_url = $2,
            product_image_url = $3,
            product_description = $4,
            updated_by_user_id = NULLIF($5, '')::uuid,
            updated_datetime = NOW()
        WHERE id = $6
        RETURNING id, product_name, product_url, product_image_url, product_description, created_datetime, updated_datetime, created_by_user_id, updated_by_user_id
    )
    SELECT 
        upd.id,
        upd.product_name,
        upd.product_url,
        upd.product_image_url,
        upd.product_description,
        upd.created_datetime,
        upd.updated_datetime,
        uc.id,
        uc.email,
        uc.password_hash,
        uu.id,
        uu.email,
        uu.password_hash
    FROM upd
    JOIN users uc ON uc.id = upd.created_by_user_id
    LEFT JOIN users uu ON uu.id = upd.updated_by_user_id
`

const SelectProductUpdatedByUserID = `
    SELECT updated_by_user_id FROM products WHERE id=$1
`

const DeleteProductByID = `
    DELETE FROM products WHERE id = $1
`
