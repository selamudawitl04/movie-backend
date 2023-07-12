
CREATE TABLE public.actors (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    image_id uuid NOT NULL
);
CREATE TABLE public.bookings (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    date timestamp with time zone DEFAULT now() NOT NULL,
    user_id uuid DEFAULT gen_random_uuid() NOT NULL,
    movie_id uuid DEFAULT gen_random_uuid() NOT NULL
);
CREATE TABLE public.contacts (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    subject text,
    body text NOT NULL,
    email text NOT NULL,
    date timestamp with time zone DEFAULT now()
);
CREATE TABLE public.directors (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    image_id uuid DEFAULT gen_random_uuid() NOT NULL
);
CREATE TABLE public.generes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL
);
CREATE TABLE public.images (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    url text DEFAULT 'default.jpg'::text NOT NULL
);
CREATE TABLE public.movies (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    duration integer NOT NULL,
    date timestamp with time zone NOT NULL,
    director_id uuid DEFAULT gen_random_uuid() NOT NULL,
    discrption text,
    cover_image uuid NOT NULL,
    status text DEFAULT 'pending'::text
);
CREATE TABLE public.movies_actors (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    movie_id uuid DEFAULT gen_random_uuid() NOT NULL,
    actor_id uuid DEFAULT gen_random_uuid() NOT NULL
);
CREATE TABLE public.movies_generes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    movie_id uuid DEFAULT gen_random_uuid() NOT NULL,
    genere_id uuid DEFAULT gen_random_uuid() NOT NULL
);
CREATE TABLE public.movies_images (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    movie_id uuid DEFAULT gen_random_uuid() NOT NULL,
    image_id uuid DEFAULT gen_random_uuid() NOT NULL
);
CREATE TABLE public.notifications (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    message text NOT NULL,
    "time" timestamp with time zone DEFAULT now() NOT NULL,
    user_id uuid DEFAULT gen_random_uuid() NOT NULL
);
CREATE TABLE public.ratings (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    movie_id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid DEFAULT gen_random_uuid() NOT NULL,
    rating numeric NOT NULL
);
CREATE TABLE public.tickets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid DEFAULT gen_random_uuid() NOT NULL,
    movie_id uuid DEFAULT gen_random_uuid() NOT NULL,
    price numeric NOT NULL,
    date timestamp with time zone DEFAULT now() NOT NULL,
    seat_number integer NOT NULL
);
CREATE TABLE public.users (
    "firstName" text NOT NULL,
    "lastName" text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    image_id uuid DEFAULT '0f6c8c0b-22fd-4593-955e-388733961620'::uuid NOT NULL,
    role text DEFAULT 'user'::text,
    "resetToken" text
);
ALTER TABLE ONLY public.actors
    ADD CONSTRAINT actors_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_pkey PRIMARY KEY (movie_id, user_id);
ALTER TABLE ONLY public.contacts
    ADD CONSTRAINT contact_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.directors
    ADD CONSTRAINT directors_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.generes
    ADD CONSTRAINT generes_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.movies_actors
    ADD CONSTRAINT movies_actors_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.movies_generes
    ADD CONSTRAINT movies_genere_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.movies_images
    ADD CONSTRAINT movies_images_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.movies
    ADD CONSTRAINT movies_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_pkey PRIMARY KEY (user_id, movie_id);
ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_pkey PRIMARY KEY (movie_id, seat_number);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT user_email_key UNIQUE (email);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pid_key UNIQUE (id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.actors
    ADD CONSTRAINT actors_image_id_fkey FOREIGN KEY (image_id) REFERENCES public.images(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON UPDATE CASCADE ON DELETE SET NULL;
ALTER TABLE ONLY public.bookings
    ADD CONSTRAINT bookings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.directors
    ADD CONSTRAINT directors_image_id_fkey FOREIGN KEY (image_id) REFERENCES public.images(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.movies_actors
    ADD CONSTRAINT movies_actors_actor_id_fkey FOREIGN KEY (actor_id) REFERENCES public.actors(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.movies_actors
    ADD CONSTRAINT movies_actors_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.movies
    ADD CONSTRAINT movies_cover_image_fkey FOREIGN KEY (cover_image) REFERENCES public.images(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.movies
    ADD CONSTRAINT movies_director_id_fkey FOREIGN KEY (director_id) REFERENCES public.directors(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE ONLY public.movies_generes
    ADD CONSTRAINT movies_genere_genere_id_fkey FOREIGN KEY (genere_id) REFERENCES public.generes(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.movies_generes
    ADD CONSTRAINT movies_genere_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.movies_images
    ADD CONSTRAINT movies_images_image_id_fkey FOREIGN KEY (image_id) REFERENCES public.images(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.movies_images
    ADD CONSTRAINT movies_images_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE SET NULL ON DELETE SET NULL;
ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON UPDATE RESTRICT ON DELETE SET NULL;
ALTER TABLE ONLY public.ratings
    ADD CONSTRAINT ratings_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE RESTRICT ON DELETE SET NULL;
ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_image_id_fkey FOREIGN KEY (image_id) REFERENCES public.images(id) ON UPDATE CASCADE ON DELETE CASCADE;
