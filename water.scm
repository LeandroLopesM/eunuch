(add-type
    "Water"
    (color 0 0 255)

    (begin) ; CreateHook
    (begin) ; DrawHook

    '(begin
        (if (and
                (cell-active? x (- y 1))
                (>=? (- y 1) 0)
            )
            (move-self x (- y 1))
            (begin)
        )
    ) ; TickHook
)