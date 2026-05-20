--
-- PostgreSQL database dump
--

\restrict 8tPKJyzO87NgJypAnBtji1ePQMvDrvKSAv6rjpl8DscPB7xGZoC5p2Abn4FEXmG

-- Dumped from database version 18.3
-- Dumped by pg_dump version 18.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: authors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.authors (
    id integer NOT NULL,
    first_name character varying NOT NULL,
    middle_name character varying,
    last_name character varying CONSTRAINT authors_surname_not_null NOT NULL,
    description text
);


ALTER TABLE public.authors OWNER TO postgres;

--
-- Name: authors_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.authors_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.authors_id_seq OWNER TO postgres;

--
-- Name: authors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.authors_id_seq OWNED BY public.authors.id;


--
-- Name: favorites; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.favorites (
    painting_id integer,
    user_id integer
);


ALTER TABLE public.favorites OWNER TO postgres;

--
-- Name: materials; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.materials (
    id integer NOT NULL,
    name character varying NOT NULL
);


ALTER TABLE public.materials OWNER TO postgres;

--
-- Name: materials_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.materials_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.materials_id_seq OWNER TO postgres;

--
-- Name: materials_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.materials_id_seq OWNED BY public.materials.id;


--
-- Name: paintings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.paintings (
    id integer NOT NULL,
    material_id integer,
    title character varying NOT NULL,
    description text,
    year integer
);


ALTER TABLE public.paintings OWNER TO postgres;

--
-- Name: paintings_authors; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.paintings_authors (
    painting_id integer,
    author_id integer
);


ALTER TABLE public.paintings_authors OWNER TO postgres;

--
-- Name: paintings_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.paintings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.paintings_id_seq OWNER TO postgres;

--
-- Name: paintings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.paintings_id_seq OWNED BY public.paintings.id;


--
-- Name: paintings_plots; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.paintings_plots (
    painting_id integer,
    plot_id integer
);


ALTER TABLE public.paintings_plots OWNER TO postgres;

--
-- Name: paintings_styles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.paintings_styles (
    painting_id integer,
    style_id integer
);


ALTER TABLE public.paintings_styles OWNER TO postgres;

--
-- Name: plots; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.plots (
    id integer NOT NULL,
    name character varying NOT NULL
);


ALTER TABLE public.plots OWNER TO postgres;

--
-- Name: plots_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.plots_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.plots_id_seq OWNER TO postgres;

--
-- Name: plots_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.plots_id_seq OWNED BY public.plots.id;


--
-- Name: styles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.styles (
    id integer NOT NULL,
    name character varying NOT NULL
);


ALTER TABLE public.styles OWNER TO postgres;

--
-- Name: styles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.styles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.styles_id_seq OWNER TO postgres;

--
-- Name: styles_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.styles_id_seq OWNED BY public.styles.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email character varying NOT NULL,
    password character varying NOT NULL,
    username character varying NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: authors id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.authors ALTER COLUMN id SET DEFAULT nextval('public.authors_id_seq'::regclass);


--
-- Name: materials id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.materials ALTER COLUMN id SET DEFAULT nextval('public.materials_id_seq'::regclass);


--
-- Name: paintings id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paintings ALTER COLUMN id SET DEFAULT nextval('public.paintings_id_seq'::regclass);


--
-- Name: plots id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plots ALTER COLUMN id SET DEFAULT nextval('public.plots_id_seq'::regclass);


--
-- Name: styles id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.styles ALTER COLUMN id SET DEFAULT nextval('public.styles_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: authors; Type: TABLE DATA; Schema: public; Owner: postgres
--

--
-- Data for Name: favorites; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: materials; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: paintings; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: paintings_authors; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: paintings_plots; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: paintings_styles; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: plots; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: styles; Type: TABLE DATA; Schema: public; Owner: postgres
--


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: authors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.authors_id_seq', 4, true);


--
-- Name: materials_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.materials_id_seq', 5, true);


--
-- Name: paintings_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.paintings_id_seq', 8, true);


--
-- Name: plots_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.plots_id_seq', 5, true);


--
-- Name: styles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.styles_id_seq', 5, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 1, false);


--
-- Name: authors authors_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.authors
    ADD CONSTRAINT authors_pkey PRIMARY KEY (id);


--
-- Name: materials materials_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.materials
    ADD CONSTRAINT materials_pkey PRIMARY KEY (id);


--
-- Name: paintings paintings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paintings
    ADD CONSTRAINT paintings_pkey PRIMARY KEY (id);


--
-- Name: plots plots_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.plots
    ADD CONSTRAINT plots_pkey PRIMARY KEY (id);


--
-- Name: styles styles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.styles
    ADD CONSTRAINT styles_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: favorites favorites_painting_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT favorites_painting_id_fkey FOREIGN KEY (painting_id) REFERENCES public.paintings(id);


--
-- Name: favorites favorites_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT favorites_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: paintings_authors paintings_authors_author_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paintings_authors
    ADD CONSTRAINT paintings_authors_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.authors(id);


--
-- Name: paintings_authors paintings_authors_painting_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paintings_authors
    ADD CONSTRAINT paintings_authors_painting_id_fkey FOREIGN KEY (painting_id) REFERENCES public.paintings(id);


--
-- Name: paintings paintings_material_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paintings
    ADD CONSTRAINT paintings_material_id_fkey FOREIGN KEY (material_id) REFERENCES public.materials(id);


--
-- Name: paintings_plots paintings_plots_painting_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paintings_plots
    ADD CONSTRAINT paintings_plots_painting_id_fkey FOREIGN KEY (painting_id) REFERENCES public.paintings(id);


--
-- Name: paintings_plots paintings_plots_plot_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paintings_plots
    ADD CONSTRAINT paintings_plots_plot_id_fkey FOREIGN KEY (plot_id) REFERENCES public.plots(id);


--
-- Name: paintings_styles paintings_styles_painting_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paintings_styles
    ADD CONSTRAINT paintings_styles_painting_id_fkey FOREIGN KEY (painting_id) REFERENCES public.paintings(id);


--
-- Name: paintings_styles paintings_styles_style_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.paintings_styles
    ADD CONSTRAINT paintings_styles_style_id_fkey FOREIGN KEY (style_id) REFERENCES public.styles(id);


--
-- PostgreSQL database dump complete
--

\unrestrict 8tPKJyzO87NgJypAnBtji1ePQMvDrvKSAv6rjpl8DscPB7xGZoC5p2Abn4FEXmG

