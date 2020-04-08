# go-registration
Registration and User 

A simple app that register users, here we implement external conection (mainly databases), an external package but still inside repo called exconn which stands for external connection and utils for a general function that could be reused.

The current main goal of this simple app is let user register with email and password, enable user to login/logout, and forgot
password serquence. We impletement basic password hasing using bcrypt and web token using JWT. We also implement multiple 
database connection, MySQL, Postgres, Mongo, and Cassandra.

More functions will be added such a function log -- any access of a function should be logged in cassandra. The objective is to gather data about input and output and execution time.

