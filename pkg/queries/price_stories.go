package queries

const ListPriceStories = `
    SELECT 
        ps.id,
        ps.created_datetime,
        ps.updated_datetime,
        p.id,
        p.product_name,
        p.product_url,
        p.product_image_url,
        p.product_description,
        p.created_datetime,
        p.updated_datetime,
        puc.id,
        puc.email,
        puc.password_hash,
        puu.id,
        puu.email,
        puu.password_hash,
        uc.id,
        uc.email,
        uc.password_hash,
        uu.id,
        uu.email,
        uu.password_hash
    FROM price_stories ps
    JOIN products p ON p.id = ps.product_id
    JOIN users uc ON uc.id = ps.created_by_user_id
    LEFT JOIN users uu ON uu.id = ps.updated_by_user_id
    JOIN users puc ON puc.id = p.created_by_user_id
    LEFT JOIN users puu ON puu.id = p.updated_by_user_id
    ORDER BY ps.created_datetime DESC
`

const GetPriceStoryByID = `
    SELECT 
        ps.id,
        ps.created_datetime,
        ps.updated_datetime,
        p.id,
        p.product_name,
        p.product_url,
        p.product_image_url,
        p.product_description,
        p.created_datetime,
        p.updated_datetime,
        puc.id,
        puc.email,
        puc.password_hash,
        puu.id,
        puu.email,
        puu.password_hash,
        uc.id,
        uc.email,
        uc.password_hash,
        uu.id,
        uu.email,
        uu.password_hash
    FROM price_stories ps
    JOIN products p ON p.id = ps.product_id
    JOIN users uc ON uc.id = ps.created_by_user_id
    LEFT JOIN users uu ON uu.id = ps.updated_by_user_id
    JOIN users puc ON puc.id = p.created_by_user_id
    LEFT JOIN users puu ON puu.id = p.updated_by_user_id
    WHERE ps.id = $1
`

const InsertPriceStoryWithUsersCTE = `
    WITH ins AS (
        INSERT INTO price_stories (product_id, created_by_user_id, updated_by_user_id)
        VALUES (
            NULLIF($1, '')::uuid,
            NULLIF($2, '')::uuid,
            NULLIF($3, '')::uuid
        )
        RETURNING id, product_id, created_by_user_id, updated_by_user_id, created_datetime, updated_datetime
    )
    SELECT 
        ins.id,
        ins.created_datetime,
        ins.updated_datetime,
        p.id,
        p.product_name,
        p.product_url,
        p.product_image_url,
        p.product_description,
        p.created_datetime,
        p.updated_datetime,
        puc.id,
        puc.email,
        puc.password_hash,
        puu.id,
        puu.email,
        puu.password_hash,
        uc.id,
        uc.email,
        uc.password_hash,
        uu.id,
        uu.email,
        uu.password_hash
    FROM ins
    JOIN products p ON p.id = ins.product_id
    JOIN users uc ON uc.id = ins.created_by_user_id
    LEFT JOIN users uu ON uu.id = ins.updated_by_user_id
    JOIN users puc ON puc.id = p.created_by_user_id
    LEFT JOIN users puu ON puu.id = p.updated_by_user_id
`

const UpdatePriceStoryWithUsersCTE = `
    WITH upd AS (
        UPDATE price_stories
        SET product_id = NULLIF($1, '')::uuid,
            updated_by_user_id = NULLIF($2, '')::uuid,
            updated_datetime = NOW()
        WHERE id = $3
        RETURNING id, product_id, created_by_user_id, updated_by_user_id, created_datetime, updated_datetime
    )
    SELECT 
        upd.id,
        upd.created_datetime,
        upd.updated_datetime,
        p.id,
        p.product_name,
        p.product_url,
        p.product_image_url,
        p.product_description,
        p.created_datetime,
        p.updated_datetime,
        puc.id,
        puc.email,
        puc.password_hash,
        puu.id,
        puu.email,
        puu.password_hash,
        uc.id,
        uc.email,
        uc.password_hash,
        uu.id,
        uu.email,
        uu.password_hash
    FROM upd
    JOIN products p ON p.id = upd.product_id
    JOIN users uc ON uc.id = upd.created_by_user_id
    LEFT JOIN users uu ON uu.id = upd.updated_by_user_id
    JOIN users puc ON puc.id = p.created_by_user_id
    LEFT JOIN users puu ON puu.id = p.updated_by_user_id
`

const DeletePriceStoryByID = `
    DELETE FROM price_stories WHERE id = $1
`
