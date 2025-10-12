package queries

const ListPricePoints = `
    SELECT 
        pp.id,
        pp.price,
        pp.timestamp,
        pp.created_datetime,
        pp.updated_datetime,
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
    FROM price_points pp
    JOIN products p ON p.id = pp.product_id
    JOIN users uc ON uc.id = pp.created_by_user_id
    LEFT JOIN users uu ON uu.id = pp.updated_by_user_id
    JOIN users puc ON puc.id = p.created_by_user_id
    LEFT JOIN users puu ON puu.id = p.updated_by_user_id
    ORDER BY pp.timestamp DESC
`

const GetPricePointByID = `
    SELECT 
        pp.id,
        pp.price,
        pp.timestamp,
        pp.created_datetime,
        pp.updated_datetime,
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
    FROM price_points pp
    JOIN products p ON p.id = pp.product_id
    JOIN users uc ON uc.id = pp.created_by_user_id
    LEFT JOIN users uu ON uu.id = pp.updated_by_user_id
    JOIN users puc ON puc.id = p.created_by_user_id
    LEFT JOIN users puu ON puu.id = p.updated_by_user_id
    WHERE pp.id = $1
`

const InsertPricePointWithUsersCTE = `
    WITH ins AS (
        INSERT INTO price_points (price, timestamp, product_id, created_by_user_id, updated_by_user_id)
        VALUES (
            $1,
            $2,
            NULLIF($3, '')::uuid,
            NULLIF($4, '')::uuid,
            NULLIF($5, '')::uuid
        )
        RETURNING id, price, timestamp, created_datetime, updated_datetime, product_id, created_by_user_id, updated_by_user_id
    )
    SELECT 
        ins.id,
        ins.price,
        ins.timestamp,
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

const UpdatePricePointWithUsersCTE = `
    WITH upd AS (
        UPDATE price_points
        SET price = $1,
            timestamp = $2,
            product_id = NULLIF($3, '')::uuid,
            updated_by_user_id = NULLIF($4, '')::uuid,
            updated_datetime = NOW()
        WHERE id = $5
        RETURNING id, price, timestamp, created_datetime, updated_datetime, product_id, created_by_user_id, updated_by_user_id
    )
    SELECT 
        upd.id,
        upd.price,
        upd.timestamp,
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

const DeletePricePointByID = `
    DELETE FROM price_points WHERE id = $1
`
