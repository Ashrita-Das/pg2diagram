%%{init: {'theme': 'dark', 'themeVariables': { 'primaryColor': '#1e1e2e', 'primaryBorderColor': '#ff5e00', 'primaryTextColor': '#cdd6f4', 'lineColor': '#ff5e00'}}}%%
erDiagram
    products ||--o{ product_categories : "category_reference_id references category_id"
    order_line_items ||--o{ customer_orders : "order_reference_id references order_id"
    order_line_items ||--o{ products : "product_reference_id references product_id"

    customer_orders {
        integer order_id
        character_varying customer_email_address
        timestamp_without_time_zone order_creation_date
    }

    order_line_items {
        integer line_item_id
        integer order_reference_id
        integer product_reference_id
        integer item_quantity
    }

    product_categories {
        integer category_id
        character_varying category_name
    }

    products {
        integer product_id
        integer category_reference_id
        character_varying product_name
        numeric product_price
    }


