# Use Server-Side Sessions

Login will use server-side sessions stored in PostgreSQL and sent to the browser through HTTP-only cookies. This keeps token handling out of the frontend, supports straightforward logout and session invalidation, and fits the locally operated deployment model better than long-lived JWTs.
