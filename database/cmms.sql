--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'SQL_ASCII';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: fuzzystrmatch; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS fuzzystrmatch WITH SCHEMA public;


--
-- Name: EXTENSION fuzzystrmatch; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION fuzzystrmatch IS 'determine similarities and distance between strings';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: component; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE component (
    machine_id integer NOT NULL,
    id integer NOT NULL,
    site_id integer NOT NULL,
    name text NOT NULL,
    descr text DEFAULT ''::text NOT NULL,
    make text DEFAULT ''::text NOT NULL,
    model text DEFAULT ''::text NOT NULL,
    serialnum text DEFAULT ''::text NOT NULL,
    picture text DEFAULT ''::text NOT NULL,
    notes text DEFAULT ''::text NOT NULL,
    qty integer DEFAULT 1 NOT NULL,
    stock_code text DEFAULT ''::text NOT NULL,
    "position" integer DEFAULT 1 NOT NULL,
    status text DEFAULT 'Running'::text NOT NULL,
    is_running boolean DEFAULT true NOT NULL,
    zindex integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.component OWNER TO postgres;

--
-- Name: component_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE component_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.component_id_seq OWNER TO postgres;

--
-- Name: component_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE component_id_seq OWNED BY component.id;


--
-- Name: component_part; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE component_part (
    component_id integer NOT NULL,
    part_id integer NOT NULL,
    qty integer DEFAULT 1 NOT NULL
);


ALTER TABLE public.component_part OWNER TO postgres;

--
-- Name: doc; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE doc (
    id integer NOT NULL,
    name text NOT NULL,
    filename text NOT NULL,
    worker boolean NOT NULL,
    sitemgr boolean NOT NULL,
    contractor boolean NOT NULL,
    type text NOT NULL,
    ref_id integer NOT NULL,
    doc_format integer DEFAULT 0 NOT NULL,
    notes text DEFAULT ''::text NOT NULL,
    filesize integer DEFAULT 0 NOT NULL,
    latest_rev integer DEFAULT 0 NOT NULL,
    created timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    user_id integer DEFAULT 0 NOT NULL,
    path text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.doc OWNER TO postgres;

--
-- Name: doc_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE doc_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.doc_id_seq OWNER TO postgres;

--
-- Name: doc_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE doc_id_seq OWNED BY doc.id;


--
-- Name: doc_rev; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE doc_rev (
    doc_id integer NOT NULL,
    id integer NOT NULL,
    revdate timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    descr text NOT NULL,
    filename text NOT NULL,
    user_id integer DEFAULT 0 NOT NULL,
    filesize integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.doc_rev OWNER TO postgres;

--
-- Name: doc_rev_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE doc_rev_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.doc_rev_id_seq OWNER TO postgres;

--
-- Name: doc_rev_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE doc_rev_id_seq OWNED BY doc_rev.id;


--
-- Name: doc_type; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE doc_type (
    id text NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.doc_type OWNER TO postgres;

--
-- Name: event; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE event (
    id integer NOT NULL,
    site_id integer NOT NULL,
    type text NOT NULL,
    machine_id integer NOT NULL,
    tool_id integer NOT NULL,
    tool_type text DEFAULT 'Tool'::text NOT NULL,
    priority integer NOT NULL,
    status text DEFAULT ''::text NOT NULL,
    startdate timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    created_by integer NOT NULL,
    allocated_by integer DEFAULT 0 NOT NULL,
    allocated_to integer DEFAULT 0 NOT NULL,
    completed timestamp with time zone,
    labour_cost numeric(12,2) DEFAULT 0.0 NOT NULL,
    material_cost numeric(12,2) DEFAULT 0.0 NOT NULL,
    other_cost numeric(12,2) DEFAULT 0.0 NOT NULL,
    notes text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.event OWNER TO postgres;

--
-- Name: event_doc; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE event_doc (
    event_id integer NOT NULL,
    doc_id integer NOT NULL,
    doc_rev_id integer NOT NULL
);


ALTER TABLE public.event_doc OWNER TO postgres;

--
-- Name: event_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE event_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.event_id_seq OWNER TO postgres;

--
-- Name: event_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE event_id_seq OWNED BY event.id;


--
-- Name: event_type; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE event_type (
    id text NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.event_type OWNER TO postgres;

--
-- Name: hashtag; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE hashtag (
    id integer NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    descr text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.hashtag OWNER TO postgres;

--
-- Name: hashtag_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE hashtag_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.hashtag_id_seq OWNER TO postgres;

--
-- Name: hashtag_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE hashtag_id_seq OWNED BY hashtag.id;


--
-- Name: machine; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE machine (
    id integer NOT NULL,
    site_id integer NOT NULL,
    name text NOT NULL,
    descr text DEFAULT ''::text NOT NULL,
    make text DEFAULT ''::text NOT NULL,
    model text DEFAULT ''::text NOT NULL,
    serialnum text NOT NULL,
    is_running boolean DEFAULT false NOT NULL,
    stopped_at timestamp with time zone,
    started_at timestamp with time zone,
    picture text DEFAULT ''::text NOT NULL,
    status text DEFAULT ''::text NOT NULL,
    notes text DEFAULT ''::text NOT NULL,
    alert_at timestamp with time zone,
    electrical text DEFAULT 'Running'::text NOT NULL,
    hydraulic text DEFAULT 'Running'::text NOT NULL,
    printer text DEFAULT 'Running'::text NOT NULL,
    console text DEFAULT 'Running'::text NOT NULL,
    rollbed text DEFAULT 'Running'::text NOT NULL,
    uncoiler text DEFAULT 'Running'::text NOT NULL,
    lube text DEFAULT 'Running'::text NOT NULL,
    tasks_to integer DEFAULT 0 NOT NULL,
    alerts_to integer DEFAULT 0 NOT NULL,
    part_class integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.machine OWNER TO postgres;

--
-- Name: machine_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE machine_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.machine_id_seq OWNER TO postgres;

--
-- Name: machine_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE machine_id_seq OWNED BY machine.id;


--
-- Name: part; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE part (
    id integer NOT NULL,
    name text NOT NULL,
    descr text DEFAULT ''::text NOT NULL,
    stock_code text NOT NULL,
    reorder_stocklevel numeric(12,2) DEFAULT 1 NOT NULL,
    reorder_qty numeric(12,2) DEFAULT 1 NOT NULL,
    latest_price numeric(12,2) DEFAULT 0 NOT NULL,
    qty_type text DEFAULT 'ea'::text NOT NULL,
    picture text DEFAULT ''::text NOT NULL,
    notes text DEFAULT ''::text NOT NULL,
    class integer DEFAULT 0 NOT NULL,
    last_price_date date,
    current_stock numeric(12,2) DEFAULT 0 NOT NULL
);


ALTER TABLE public.part OWNER TO postgres;

--
-- Name: part_class; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE part_class (
    id integer NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    descr text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.part_class OWNER TO postgres;

--
-- Name: part_class_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE part_class_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.part_class_id_seq OWNER TO postgres;

--
-- Name: part_class_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE part_class_id_seq OWNED BY part_class.id;


--
-- Name: part_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE part_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.part_id_seq OWNER TO postgres;

--
-- Name: part_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE part_id_seq OWNED BY part.id;


--
-- Name: part_price; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE part_price (
    part_id integer NOT NULL,
    datefrom timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    price numeric(12,2) DEFAULT 0 NOT NULL
);


ALTER TABLE public.part_price OWNER TO postgres;

--
-- Name: part_stock; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE part_stock (
    part_id integer NOT NULL,
    datefrom timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    stock_level numeric(12,2) DEFAULT 0 NOT NULL
);


ALTER TABLE public.part_stock OWNER TO postgres;

--
-- Name: part_vendor; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE part_vendor (
    part_id integer NOT NULL,
    vendor_id integer NOT NULL,
    vendor_code text NOT NULL,
    latest_price numeric(12,2) NOT NULL
);


ALTER TABLE public.part_vendor OWNER TO postgres;

--
-- Name: sched_control; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sched_control (
    id integer NOT NULL,
    last_run date
);


ALTER TABLE public.sched_control OWNER TO postgres;

--
-- Name: sched_control_task; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sched_control_task (
    task_id integer NOT NULL,
    last_gen date,
    last_jobcount integer
);


ALTER TABLE public.sched_control_task OWNER TO postgres;

--
-- Name: sched_task; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sched_task (
    id integer NOT NULL,
    machine_id integer NOT NULL,
    comp_type text DEFAULT 'C'::text NOT NULL,
    tool_id integer NOT NULL,
    component text DEFAULT ''::text NOT NULL,
    descr text DEFAULT ''::text NOT NULL,
    startdate date,
    oneoffdate date,
    freq text DEFAULT 'R'::text NOT NULL,
    parent_task integer,
    days integer,
    count integer,
    week integer,
    duration_days integer DEFAULT 1 NOT NULL,
    labour_cost numeric(12,2) NOT NULL,
    material_cost numeric(12,2) NOT NULL,
    other_cost_desc text[],
    other_cost numeric(12,2)[],
    last_generated date,
    weekday integer,
    user_id integer DEFAULT 0 NOT NULL,
    paused boolean DEFAULT true NOT NULL
);


ALTER TABLE public.sched_task OWNER TO postgres;

--
-- Name: sched_task_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE sched_task_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sched_task_id_seq OWNER TO postgres;

--
-- Name: sched_task_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE sched_task_id_seq OWNED BY sched_task.id;


--
-- Name: sched_task_part; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sched_task_part (
    task_id integer NOT NULL,
    part_id integer NOT NULL,
    qty numeric(12,2) NOT NULL,
    notes text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.sched_task_part OWNER TO postgres;

--
-- Name: site; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE site (
    id integer NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    address text DEFAULT ''::text NOT NULL,
    phone text DEFAULT ''::text NOT NULL,
    fax text DEFAULT ''::text NOT NULL,
    image text DEFAULT ''::text NOT NULL,
    parent_site integer DEFAULT 0 NOT NULL,
    notes text DEFAULT ''::text NOT NULL,
    stock_site integer DEFAULT 0 NOT NULL,
    x integer DEFAULT 0 NOT NULL,
    y integer DEFAULT 0 NOT NULL,
    alerts_to integer DEFAULT 0 NOT NULL,
    tasks_to integer DEFAULT 0 NOT NULL,
    manager integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.site OWNER TO postgres;

--
-- Name: site_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE site_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.site_id_seq OWNER TO postgres;

--
-- Name: site_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE site_id_seq OWNED BY site.id;


--
-- Name: site_layout; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE site_layout (
    site_id integer NOT NULL,
    seq integer NOT NULL,
    machine_id integer NOT NULL,
    span integer NOT NULL
);


ALTER TABLE public.site_layout OWNER TO postgres;

--
-- Name: skill; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE skill (
    id integer NOT NULL,
    name text NOT NULL,
    notes text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.skill OWNER TO postgres;

--
-- Name: skill_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE skill_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.skill_id_seq OWNER TO postgres;

--
-- Name: skill_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE skill_id_seq OWNED BY skill.id;


--
-- Name: sm_component; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sm_component (
    task_id integer NOT NULL,
    machine_id integer NOT NULL,
    component text NOT NULL,
    completed timestamp with time zone,
    mins_spent integer DEFAULT 0 NOT NULL,
    labour_cost numeric(12,2) NOT NULL,
    materal_cost numeric(12,2) NOT NULL,
    notes text
);


ALTER TABLE public.sm_component OWNER TO postgres;

--
-- Name: sm_component_item; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sm_component_item (
    task_id integer NOT NULL,
    component text NOT NULL,
    seq integer NOT NULL,
    notes text,
    done boolean DEFAULT false NOT NULL
);


ALTER TABLE public.sm_component_item OWNER TO postgres;

--
-- Name: sm_machine; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sm_machine (
    task_id integer NOT NULL,
    machine_id integer NOT NULL,
    completed timestamp with time zone
);


ALTER TABLE public.sm_machine OWNER TO postgres;

--
-- Name: sm_machine_item; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sm_machine_item (
    task_id integer NOT NULL,
    machine_id integer NOT NULL,
    seq integer NOT NULL,
    notes text,
    done boolean DEFAULT false NOT NULL
);


ALTER TABLE public.sm_machine_item OWNER TO postgres;

--
-- Name: sm_parts; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sm_parts (
    task_id integer NOT NULL,
    part_id integer NOT NULL,
    date timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    qty numeric(12,2) NOT NULL,
    value numeric(12,2) NOT NULL
);


ALTER TABLE public.sm_parts OWNER TO postgres;

--
-- Name: sm_task; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sm_task (
    id integer NOT NULL,
    user_id integer NOT NULL,
    date timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    type text DEFAULT ''::text NOT NULL,
    week integer DEFAULT 1 NOT NULL,
    completed timestamp with time zone,
    escalate_date timestamp with time zone,
    escalate_user integer,
    status text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.sm_task OWNER TO postgres;

--
-- Name: sm_task_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE sm_task_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sm_task_id_seq OWNER TO postgres;

--
-- Name: sm_task_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE sm_task_id_seq OWNED BY sm_task.id;


--
-- Name: sm_tool; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sm_tool (
    task_id integer NOT NULL,
    machine_id integer NOT NULL,
    tool_id integer NOT NULL,
    completed timestamp with time zone,
    mins_spent integer DEFAULT 0 NOT NULL,
    labour_cost numeric(12,2) NOT NULL,
    materal_cost numeric(12,2) NOT NULL,
    notes text
);


ALTER TABLE public.sm_tool OWNER TO postgres;

--
-- Name: sm_tool_item; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sm_tool_item (
    task_id integer NOT NULL,
    tool_id integer NOT NULL,
    seq integer NOT NULL,
    notes text,
    done boolean DEFAULT false NOT NULL
);


ALTER TABLE public.sm_tool_item OWNER TO postgres;

--
-- Name: stock_level; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE stock_level (
    part_id integer NOT NULL,
    site_id integer NOT NULL,
    datefrom timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    qty numeric(12,2) NOT NULL,
    notes text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.stock_level OWNER TO postgres;

--
-- Name: stock_level_part_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE stock_level_part_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.stock_level_part_id_seq OWNER TO postgres;

--
-- Name: stock_level_part_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE stock_level_part_id_seq OWNED BY stock_level.part_id;


--
-- Name: sys_log; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE sys_log (
    id integer NOT NULL,
    status integer DEFAULT 0 NOT NULL,
    type text NOT NULL,
    ref_type character(1) NOT NULL,
    ref_id integer NOT NULL,
    logdate timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    ip text NOT NULL,
    descr text NOT NULL,
    user_id integer NOT NULL,
    username text DEFAULT ''::text NOT NULL,
    before text DEFAULT ''::text NOT NULL,
    after text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.sys_log OWNER TO postgres;

--
-- Name: sys_log_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE sys_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sys_log_id_seq OWNER TO postgres;

--
-- Name: sys_log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE sys_log_id_seq OWNED BY sys_log.id;


--
-- Name: task; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE task (
    id integer NOT NULL,
    sched_id integer DEFAULT 0 NOT NULL,
    machine_id integer NOT NULL,
    tool_id integer NOT NULL,
    comp_type text DEFAULT 'C'::text NOT NULL,
    component text DEFAULT ''::text NOT NULL,
    descr text DEFAULT ''::text NOT NULL,
    log text DEFAULT ''::text NOT NULL,
    created_date timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    startdate date,
    due_date date,
    escalate_date date,
    assigned_by integer,
    assigned_to integer,
    assigned_date timestamp with time zone,
    completed_date timestamp with time zone,
    has_issue boolean DEFAULT false NOT NULL,
    issue_resolved_date timestamp with time zone,
    labour_est numeric(12,2) DEFAULT 0 NOT NULL,
    material_est numeric(12,2) DEFAULT 0 NOT NULL,
    labour_cost numeric(12,2) DEFAULT 0 NOT NULL,
    material_cost numeric(12,2) DEFAULT 0 NOT NULL,
    other_cost_desc text[],
    other_cost numeric(12,2)[]
);


ALTER TABLE public.task OWNER TO postgres;

--
-- Name: task_check; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE task_check (
    task_id integer NOT NULL,
    seq integer NOT NULL,
    descr text NOT NULL,
    done boolean DEFAULT false NOT NULL,
    done_date timestamp with time zone
);


ALTER TABLE public.task_check OWNER TO postgres;

--
-- Name: task_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE task_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.task_id_seq OWNER TO postgres;

--
-- Name: task_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE task_id_seq OWNED BY task.id;


--
-- Name: task_part; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE task_part (
    task_id integer NOT NULL,
    part_id integer NOT NULL,
    qty numeric(12,2) NOT NULL,
    notes text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.task_part OWNER TO postgres;

--
-- Name: user_log; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_log (
    id integer NOT NULL,
    logdate timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    duration text,
    ms integer,
    func text,
    input text,
    output text
);


ALTER TABLE public.user_log OWNER TO postgres;

--
-- Name: user_log_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE user_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_log_id_seq OWNER TO postgres;

--
-- Name: user_log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE user_log_id_seq OWNED BY user_log.id;


--
-- Name: user_role; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_role (
    user_id integer NOT NULL,
    site_id integer NOT NULL,
    worker boolean DEFAULT false NOT NULL,
    sitemgr boolean NOT NULL,
    contractor boolean NOT NULL
);


ALTER TABLE public.user_role OWNER TO postgres;

--
-- Name: user_site; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_site (
    user_id integer NOT NULL,
    site_id integer NOT NULL,
    role text NOT NULL
);


ALTER TABLE public.user_site OWNER TO postgres;

--
-- Name: user_skill; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE user_skill (
    user_id integer NOT NULL,
    skill_id integer NOT NULL
);


ALTER TABLE public.user_skill OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE users (
    id integer NOT NULL,
    username character varying(32) NOT NULL,
    passwd character varying(32) NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    address text DEFAULT ''::text NOT NULL,
    email text DEFAULT ''::text NOT NULL,
    sms text DEFAULT ''::text NOT NULL,
    site_id integer DEFAULT 0 NOT NULL,
    role text DEFAULT 'Public'::text NOT NULL,
    notes text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: vendor; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE vendor (
    id integer NOT NULL,
    name text NOT NULL,
    descr text DEFAULT ''::text NOT NULL,
    address text DEFAULT ''::text NOT NULL,
    phone text DEFAULT ''::text NOT NULL,
    fax text DEFAULT ''::text NOT NULL,
    contact_name text DEFAULT ''::text NOT NULL,
    contact_email text DEFAULT ''::text NOT NULL,
    orders_email text DEFAULT ''::text NOT NULL,
    rating text DEFAULT ''::text NOT NULL,
    notes text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.vendor OWNER TO postgres;

--
-- Name: vendor_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE vendor_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.vendor_id_seq OWNER TO postgres;

--
-- Name: vendor_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE vendor_id_seq OWNED BY vendor.id;


--
-- Name: vendor_price; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE vendor_price (
    part_id integer NOT NULL,
    vendor_id integer NOT NULL,
    datefrom timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    price numeric(12,2) NOT NULL,
    min_qty numeric(12,2) NOT NULL,
    notes text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.vendor_price OWNER TO postgres;

--
-- Name: wo_assignee; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE wo_assignee (
    id integer NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.wo_assignee OWNER TO postgres;

--
-- Name: wo_docs; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE wo_docs (
    id integer NOT NULL,
    doc_id integer NOT NULL
);


ALTER TABLE public.wo_docs OWNER TO postgres;

--
-- Name: wo_skills; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE wo_skills (
    id integer NOT NULL,
    skill_id integer NOT NULL
);


ALTER TABLE public.wo_skills OWNER TO postgres;

--
-- Name: workorder; Type: TABLE; Schema: public; Owner: postgres; Tablespace: 
--

CREATE TABLE workorder (
    id integer NOT NULL,
    event_id integer DEFAULT 0 NOT NULL,
    startdate timestamp with time zone DEFAULT ('now'::text)::timestamp without time zone NOT NULL,
    est_duration integer DEFAULT 0 NOT NULL,
    actual_duration integer DEFAULT 0 NOT NULL,
    descr text DEFAULT ''::text NOT NULL,
    status text DEFAULT ''::text NOT NULL,
    notes text DEFAULT ''::text NOT NULL
);


ALTER TABLE public.workorder OWNER TO postgres;

--
-- Name: workorder_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE workorder_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.workorder_id_seq OWNER TO postgres;

--
-- Name: workorder_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE workorder_id_seq OWNED BY workorder.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY component ALTER COLUMN id SET DEFAULT nextval('component_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY doc ALTER COLUMN id SET DEFAULT nextval('doc_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY doc_rev ALTER COLUMN id SET DEFAULT nextval('doc_rev_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY event ALTER COLUMN id SET DEFAULT nextval('event_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY hashtag ALTER COLUMN id SET DEFAULT nextval('hashtag_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY machine ALTER COLUMN id SET DEFAULT nextval('machine_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY part ALTER COLUMN id SET DEFAULT nextval('part_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY part_class ALTER COLUMN id SET DEFAULT nextval('part_class_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY sched_task ALTER COLUMN id SET DEFAULT nextval('sched_task_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY site ALTER COLUMN id SET DEFAULT nextval('site_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY skill ALTER COLUMN id SET DEFAULT nextval('skill_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY sm_task ALTER COLUMN id SET DEFAULT nextval('sm_task_id_seq'::regclass);


--
-- Name: part_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY stock_level ALTER COLUMN part_id SET DEFAULT nextval('stock_level_part_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY sys_log ALTER COLUMN id SET DEFAULT nextval('sys_log_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY task ALTER COLUMN id SET DEFAULT nextval('task_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY user_log ALTER COLUMN id SET DEFAULT nextval('user_log_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY vendor ALTER COLUMN id SET DEFAULT nextval('vendor_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY workorder ALTER COLUMN id SET DEFAULT nextval('workorder_id_seq'::regclass);


--
-- Data for Name: component; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY component (machine_id, id, site_id, name, descr, make, model, serialnum, picture, notes, qty, stock_code, "position", status, is_running, zindex) FROM stdin;
26	242	8	Service Hole #2	Service Hole #2						1	SBSPL92-01-0000	4	Running	t	0
26	243	8	Curl #2	Curl #2						1	SBSSD75-01-0001	11	Running	t	0
29	234	6	Down Dimple - Fixed	Down Dimple - Fixed						1	SBSCH-01-0001	1	Running	t	0
30	235	6	Large Crush	Large Crush						1	SBSWB-02-0001	2	Running	t	0
5	14	2	Single Dimple Square	Single Dimple Square						1	SBSPL75-01-0000 	1	Running	t	0
6	21	2	Service Hole #1	Service Hole #1						1	SBSSD75-01-0000 	2	Running	t	0
6	236	2	Brick Tie	Brick Tie						1	SBSSD75-01-0001	1	Running	t	0
19	80	7	Quad Dimple	Quad Dimple						1	SBSSD75-02-0000 	3	Running	t	0
23	246	8	Brick Tie	Brick Tie						1	SBSSD75-01-0001	1	Running	t	0
23	244	8	Curl #2	Curl #2						1	SBSSD75-01-0001	8	Running	t	0
23	245	8	Service Hole #2	Service Hole #2						1	SBSSD75-01-0001	4	Running	t	0
28	247	6	Brick Tie	Brick Tie						1	SBSSD75-01-0001	1	Running	t	0
28	249	6	Curl #2	Curl #2						1	SBSSD75-01-0001	8	Running	t	0
28	248	6	Service Hole #2	Service Hole #2						1	SBSSD75-01-0001	4	Running	t	0
18	78	7	Guillo	Guillo						1	SBSPL75-07-0000	7	Running	t	0
18	72	7	 Single Dimple Square	 Single Dimple Square						1	SBSPL75-01-0000 	1	Running	t	0
6	238	2	Curl #2	Curl #2						1	SBSSD75-01-0001	8	Running	t	0
19	239	7	Brick Tie	Brick Tie						1	SBSSD75-01-0000	1	Running	t	0
19	240	7	Curl #2	Curl #2						1	SBSSD75-01-0000	8	Running	t	0
19	241	7	Service Hole #2	Service Hole #2						1	SBSSD75-01-0000	4	Running	t	0
37	190	9	H-Cut	H-Cut						1	SBSFF250-01-000	1	Running	t	0
38	200	9	 Single Dimple Square	 Single Dimple Square						1	SBSPL75-01-0000 	1	Running	t	0
39	207	9	Service Hole	Service Hole						1	SBSSD75-01-0000 	1	Running	t	0
42	232	9	Guillo	Guillo						1	SBSPM-01-000	1	Running	t	0
43	233	9	Guillo	Guillo						1	SBSTH22-01-0000	1	Running	t	0
37	195	9	Service Hole	Service Hole						1	SBSFF250-12-000	6	Running	t	0
37	196	9	Down Dimple	Down Dimple						1	SBSFF250-15-000	7	Running	t	0
37	197	9	Up Dimple	Up Dimple						1	SBSFF250-17-000	8	Running	t	0
37	198	9	Guillo Joist	Guillo Joist						1	SBSFF250-19-000	9	Running	t	0
5	17	2	Tie Down Slot	Tie Down Slot						1	SBSPL75-04-0000 	3	Running	t	0
6	23	2	Single Dimple & Rib	Single Dimple & Rib						1	SBSSD75-03-0000 	5	Running	t	0
6	24	2	Curl #1	Curl #1						2	SBSSD75-04-0000 	6	Running	t	0
6	25	2	Guillo	Guillo						1	SBSSD75-05-0000	7	Running	t	0
6	22	2	Quad Dimple	Quad Dimple						1	SBSSD75-02-0000 	3	Running	t	0
19	79	7	Service Hole #1	Service Hole #1						1	SBSSD75-01-0000 	2	Running	t	0
14	62	4	Tie Down Slot	Tie Down Slot						1	SBSCH-03-0000	3	Needs Attention	t	0
5	20	2	Guillo	Guillo						1	SBSPL75-07-0000	7	Running	t	0
14	61	4	Half Notch	Half Notch						1	SBSCH-02-0000	2	Needs Attention	t	0
14	63	4	Up Dimple	Up Dimple						1	SBSCH-04-0000	4	Needs Attention	t	0
21	99	7	Pierce Location 	Pierce Location 						1	SBSWB-01-0000	1	Running	t	0
21	100	7	Crush	Crush						1	SBSWB-02-0000	2	Running	t	0
21	101	7	Guillo	Guillo						1	SBSWB-03-0000	3	Running	t	0
25	131	8	Pierce Location 	Pierce Location 						1	SBSWB-01-0000	1	Running	t	0
26	136	8	Service Hole #1	Service Hole #1						2	SBSSD92-01-0000	2	Running	t	0
11	49	2	Guillo Lintel	Guillo Lintel						1	SBSVL-01-0221	2	Running	t	0
13	52	4	Quad Dimple	Quad Dimple						1	SBSSD92-02-0000	2	Running	t	0
13	53	4	Single Dimple Crush	Single Dimple Crush						1	SBSSD92-03-0000	3	Running	t	0
13	54	4	Curl	Curl						2	SBSSD92-04-0000	4	Running	t	0
13	56	4	Single Dimple Square	Single Dimple Square						1	SBSPL92-01-0000	6	Running	t	0
13	58	4	Tie Down Slot	Tie Down Slot						1	SBSPL92-04-0000	8	Running	t	0
13	59	4	Nogging	Nogging						1	SBSPL92-06-0000	9	Running	t	0
14	60	4	Down Dimple	Down Dimple						1	SBSCH-01-0000	1	Running	t	0
14	64	4	Full Notch	Full Notch						1	SBSCH-05-0000	5	Running	t	0
14	65	4	Right Angle Guillo	Right Angle Guillo						1	SBSCH-06-0000	6	Running	t	0
14	66	4	Straight Guillo	Straight Guillo						1	SBSCH-07-0000	7	Running	t	0
14	67	4	Left Angle Guillo	Left Angle Guillo						1	SBSCH-08-0000	8	Running	t	0
18	76	7	Curl	Curl						1	SBSPL75-05-0000	5	Running	t	0
13	51	4	Service Hole	Service Hole						2	SBSSD92-01-0000	1	Needs Attention	t	0
10	47	2	Guillo	Guillo						1	SBSTH22-01-0000	1	Needs Attention	t	0
13	55	4	Guillo	Guillo						1	SBSSD92-05-0000	5	Needs Attention	t	0
5	15	2	Service Hole / Curl	Service Hole / Curl						1	SBSPL75-02-0000 	2	Needs Attention	t	0
6	237	2	Service Hole #2	Service Hole #2						1	SBSSD75-01-0001	4	Needs Attention	t	0
11	48	2	Guillo Valley	Guillo Valley						1	SBSVL-01-0201	1	Needs Attention	t	0
13	57	4	Notch	Notch						1	SBSPL92-03-0000	7	Needs Attention	t	0
18	77	7	Nogging	Nogging						1	SBSPL75-06-0000 	6	Running	t	0
18	73	7	Service Hole	Service Hole						1	SBSPL75-02-0000 	2	Needs Attention	t	0
5	16	2	Notch	Notch						1	SBSPL75-03-0000 	4	Needs Attention	t	0
18	74	7	Notch	Notch						1	SBSPL75-03-0000 	3	Needs Attention	t	0
5	19	2	Nogging	Nogging						1	SBSPL75-06-0000 	6	Needs Attention	t	0
24	122	8	Straight Guillo	Straight Guillo						1	SBSCH-07-0000	7	Maintenance Pending	t	0
25	132	8	Crush	Crush						1	SBSWB-02-0000	2	Running	t	0
22	104	8	 Single Dimple Square	 Single Dimple Square						1	SBSPL75-01-0000 	1	Running	t	0
29	159	6	Tie Down Slot	Tie Down Slot						1	SBSCH-03-0000	3	Needs Attention	t	0
25	133	8	Guillo	Guillo						1	SBSWB-03-0000	3	Running	t	0
4	5	2	Fold Rib	Fold Rib						1	SBSFF150-06-000 	5	Running	t	0
24	119	8	Up Dimple	Up Dimple						1	SBSCH-04-0000	3	Running	t	0
23	111	8	Service Hole #1	Service Hole #1						1	SBSSD75-01-0000 	2	Running	t	0
23	112	8	Quad Dimple	Quad Dimple						1	SBSSD75-02-0000 	3	Running	t	0
23	113	8	Single Dimple & Rib	Single Dimple & Rib						1	SBSSD75-03-0000 	5	Running	t	0
22	105	8	Service Hole	Service Hole						1	SBSPL75-02-0000 	2	Running	t	0
22	106	8	Notch	Notch						1	SBSPL75-03-0000 	3	Running	t	0
22	107	8	Tie Down Slot	Tie Down Slot						1	SBSPL75-04-0000 	4	Running	t	0
22	108	8	Curl	Curl						1	SBSPL75-05-0000	5	Running	t	0
22	109	8	Nogging	Nogging						1	SBSPL75-06-0000 	6	Running	t	0
22	110	8	Guillo	Guillo						1	SBSPL75-07-0000	7	Running	t	0
23	114	8	Curl #1	Curl #1						2	SBSSD75-04-0000 	6	Running	t	0
23	115	8	Guillo	Guillo						1	SBSSD75-05-0000	7	Running	t	0
24	120	8	Full Notch	Full Notch						1	SBSCH-05-0000	5	Running	t	0
29	158	6	Half Notch	Half Notch					<b>Maintenance Plan for the Half Notch Tool</b><br><br><ol><li>Prepare work area</li><li>Do the work</li><li>plug the machine back in</li><li>Done !</li></ol>	1	SBSCH-02-0000	2	Running	t	0
24	123	8	Left Angle Guillo	Left Angle Guillo						1	SBSCH-08-0000	8	Running	t	0
27	145	6	 Single Dimple Square	 Single Dimple Square						1	SBSPL75-01-0000 	1	Running	t	0
27	147	6	Notch	Notch						1	SBSPL75-03-0000 	3	Running	t	0
8	41	2	Pierce Location 	Pierce Location 						1	SBSWB-01-0000	1	Running	t	0
8	42	2	Crush	Crush						1	SBSWB-02-0000	2	Running	t	0
8	43	2	Guillo	Guillo						1	SBSWB-03-0000	3	Running	t	0
29	157	6	Down Dimple - Moveable	Down Dimple - Moveable						3	SBSCH-01-0000	1	Running	t	0
19	83	7	Guillo	Guillo						1	SBSSD75-05-0000	7	Running	t	0
27	148	6	Tie Down Slot	Tie Down Slot						1	SBSPL75-04-0000 	4	Running	t	0
27	149	6	Curl	Curl						1	SBSPL75-05-0000	5	Running	t	0
27	150	6	Nogging	Nogging						1	SBSPL75-06-0000 	6	Running	t	0
27	151	6	Guillo	Guillo						1	SBSPL75-07-0000	7	Running	t	0
36	179	9	Pier Slot	Pier Slot						2	SBSFF150-04-000 	3	Running	t	0
36	180	9	Tie Down Slot	Tie Down Slot						1	SBSFF150-23-000 	4	Running	t	0
36	181	9	Fold Rib	Fold Rib						1	SBSFF150-06-000 	5	Running	t	0
36	182	9	Dimple Joist	Dimple Joist						1	SBSFF150-08-000	6	Running	t	0
36	184	9	Service Hole	Service Hole						1	SBSFF150-12-000 	8	Running	t	0
36	185	9	Swage	Swage						1	SBSFF150-13-000 	9	Running	t	0
36	186	9	Up Dimple	Up Dimple						1	SBSFF150-17-000 	10	Running	t	0
36	187	9	Down Dimple	Down Dimple						1	SBSFF150-15-000 	11	Running	t	0
36	188	9	Guillo Joist	Guillo Joist						1	SBSFF150-19-000 	12	Running	t	0
36	189	9	Guillo Bearer	Guillo Bearer						1	SBSFF150-20-000 	13	Running	t	0
37	191	9	Notch	Notch						1	SBSFF250-02-000	2	Running	t	0
37	192	9	Fold Rib	Fold Rib						1	SBSFF250-06-000	3	Running	t	0
37	193	9	Dimple Joist 	Dimple Joist 						1	SBSFF250-08-000	4	Running	t	0
37	194	9	Dimple Bearer	Dimple Bearer						1	SBSFF250-10-000	5	Running	t	0
36	177	9	H-Cut	H-Cut						1	SBSFF150-01-000	1	Running	t	0
36	183	9	Dimple Bearer	Dimple Bearer						1	SBSFF150-10-000 	7	Maintenance Pending	t	0
41	227	9	Pierce Location 	Pierce Location 						1	SBSWB-01-0000	2	Running	t	0
41	228	9	Crush	Crush						1	SBSWB-02-0000	3	Running	t	0
41	229	9	Guillo	Guillo						1	SBSWB-03-0000	4	Running	t	0
26	141	8	Single Dimple Square	Single Dimple Square						1	SBSPL92-01-0000	1	Running	t	0
37	199	9	Guillo Bearer	Guillo Bearer						1	SBSFF250-20-000	10	Running	t	0
38	201	9	Service Hole	Service Hole						1	SBSPL75-02-0000 	2	Running	t	0
38	202	9	Notch	Notch						1	SBSPL75-03-0000 	3	Running	t	0
38	203	9	Tie Down Slot	Tie Down Slot						1	SBSPL75-04-0000 	4	Running	t	0
38	204	9	Curl	Curl						1	SBSPL75-05-0000	5	Running	t	0
38	205	9	Nogging	Nogging						1	SBSPL75-06-0000 	6	Running	t	0
38	206	9	Guillo	Guillo						1	SBSPL75-07-0000	7	Running	t	0
39	209	9	Single Dimple & Rib	Single Dimple & Rib						1	SBSSD75-03-0000 	3	Running	t	0
39	210	9	Curl	Curl						2	SBSSD75-04-0000 	4	Running	t	0
39	211	9	Guillo	Guillo						1	SBSSD75-05-0000	5	Running	t	0
26	137	8	Quad Dimple	Quad Dimple						1	SBSSD92-02-0000	3	Running	t	0
24	121	8	Right Angle Guillo	Right Angle Guillo						1	SBSCH-06-0000	6	Running	t	0
27	146	6	Service Hole	Service Hole						1	SBSPL75-02-0000 	2	Running	t	0
26	138	8	Single Dimple Crush	Single Dimple Crush						1	SBSSD92-03-0000	7	Running	t	0
26	139	8	Curl #1	Curl #1						2	SBSSD92-04-0000	9	Running	t	0
30	174	6	Guillo	Guillo						1	SBSWB-03-0000	4	Running	t	0
19	81	7	Single Dimple & Rib	Single Dimple & Rib						1	SBSSD75-03-0000 	5	Needs Attention	t	0
26	140	8	Guillo	Guillo						1	SBSSD92-05-0000	10	Running	t	0
26	142	8	Notch	Notch						1	SBSPL92-03-0000	6	Running	t	0
29	161	6	Full Notch	Full Notch						1	SBSCH-05-0000	5	Running	t	0
29	163	6	Straight Guillo	Straight Guillo						1	SBSCH-07-0000	7	Running	t	0
29	164	6	Left Angle Guillo	Left Angle Guillo						1	SBSCH-08-0000	8	Running	t	0
36	178	9	Notch	Notch						1	SBSFF150-02-000 	2	Running	t	0
39	208	9	Quad Dimple	Quad Dimple						1	SBSSD75-02-0000 	2	Running	t	0
40	213	9	Half Notch	Half Notch						1	SBSCH-02-0000	2	Running	t	0
40	214	9	Tie Down Slot	Tie Down Slot						1	SBSCH-03-0000	3	Running	t	0
40	215	9	Up Dimple	Up Dimple						1	SBSCH-04-0000	4	Running	t	0
40	216	9	Full Notch	Full Notch						1	SBSCH-05-0000	5	Running	t	0
40	217	9	Right Angle Guillo	Right Angle Guillo						1	SBSCH-06-0000	6	Running	t	0
40	218	9	Straight Guillo	Straight Guillo						1	SBSCH-07-0000	7	Running	t	0
40	219	9	Left Angle Guillo	Left Angle Guillo						1	SBSCH-08-0000	8	Running	t	0
26	143	8	Tie Down Slot	Tie Down Slot						1	SBSPL92-04-0000	5	Running	t	0
26	144	8	Nogging	Nogging						1	SBSPL92-06-0000	8	Running	t	0
30	172	6	Pierce Location 	Pierce Location 						1	SBSWB-01-0000	1	Running	t	0
28	152	6	Service Hole #1	Service Hole #1						1	SBSSD75-01-0000 	2	Running	t	0
30	173	6	Small Crush	Small Crush						1	SBSWB-02-0000	3	Running	t	0
28	153	6	Quad Dimple	Quad Dimple						1	SBSSD75-02-0000 	3	Running	t	0
28	154	6	Single Dimple & Rib	Single Dimple & Rib						1	SBSSD75-03-0000 	5	Running	t	0
28	155	6	Curl #1	Curl #1						2	SBSSD75-04-0000 	6	Running	t	0
28	156	6	Guillo	Guillo						1	SBSSD75-05-0000	7	Running	t	0
24	117	8	Half Notch	Half Notch						1	SBSCH-02-0000	4	Running	t	0
24	118	8	Tie Down Slot	Tie Down Slot						1	SBSCH-03-0000	2	Running	t	0
29	162	6	Right Angle Guillo	Right Angle Guillo						1	SBSCH-06-0000	6	Running	t	0
4	13	2	Guillo Bearer	Guillo Bearer						1	SBSFF150-20-000 	12	Running	t	0
20	84	7	Down Dimple	Down Dimple						1	SBSCH-01-0000	1	Running	t	0
20	89	7	Right Angle Guillo	Right Angle Guillo						1	SBSCH-06-0000	6	Running	t	0
40	212	9	Down Dimple	Down Dimple						3	SBSCH-01-0000	1	Stopped	f	0
24	116	8	Down Dimple	Down Dimple						1	SBSCH-01-0000	1	Running	t	0
7	30	2	Full Notch	Full Notch						1	SBSCH-05-0000	5	Running	t	0
7	31	2	Right Angle Guillo	Right Angle Guillo						1	SBSCH-06-0000	6	Running	t	0
7	26	2	Down Dimple	Down Dimple						3	SBSCH-01-0000	1	Running	t	0
4	10	2	Up Dimple	Up Dimple						1	SBSFF150-17-000 	10	Running	t	0
4	6	2	Dimple Joist	Dimple Joist						1	SBSFF150-08-000	6	Running	t	0
4	8	2	Service Hole	Service Hole						1	SBSFF150-12-000 	8	Running	t	0
4	4	2	Tie Down Slot	Tie Down Slot						1	SBSFF150-23-000 	4	Running	t	0
4	11	2	Down Dimple	Down Dimple						1	SBSFF150-15-000 	11	Running	t	0
4	1	2	H-Cut	H-Cut						1	SBSFF150-01-000	1	Running	t	0
4	3	2	Pier Slot	Pier Slot						2	SBSFF150-04-000 	3	Running	t	0
4	7	2	Dimple Bearer	Dimple Bearer						1	SBSFF150-10-000 	6	Running	t	0
4	12	2	Guillo Joist	Guillo Joist						1	SBSFF150-19-000 	12	Running	t	0
7	28	2	Tie Down Slot	Tie Down Slot						1	SBSCH-03-0000	2	Running	t	0
7	29	2	Up Dimple	Up Dimple						1	SBSCH-04-0000	3	Running	t	0
7	27	2	Half Notch	Half Notch						1	SBSCH-02-0000	4	Running	t	0
7	33	2	Left Angle Guillo	Left Angle Guillo						1	SBSCH-08-0000	8	Running	t	0
7	32	2	Straight Guillo	Straight Guillo						1	SBSCH-07-0000	7	Running	t	0
20	85	7	Half Notch	Half Notch						1	SBSCH-02-0000	4	Running	t	0
20	86	7	Tie Down Slot	Tie Down Slot						1	SBSCH-03-0000	2	Needs Attention	t	0
20	87	7	Up Dimple	Up Dimple						1	SBSCH-04-0000	3	Running	t	0
20	91	7	Left Angle Guillo	Left Angle Guillo						1	SBSCH-08-0000	8	Running	t	0
20	90	7	Straight Guillo	Straight Guillo						1	SBSCH-07-0000	7	Running	t	0
15	69	4	Pierce Location 	Pierce Location 						1	SBSWB-01-0000	2	Running	t	0
15	70	4	Crush	Crush						1	SBSWB-02-0000	3	Running	t	0
15	71	4	Guillo	Guillo						1	SBSWB-03-0000	4	Running	t	0
4	2	2	Notch	Notch						1	SBSFF150-02-000 	2	Needs Attention	t	0
9	46	2	Guillo	Guillo						1	SBSPM-01-000	1	Needs Attention	t	0
4	9	2	Swage	Swage						1	SBSFF150-13-000 	9	Needs Attention	t	0
29	160	6	Up Dimple	Up Dimple						1	SBSCH-04-0000	4	Needs Attention	t	0
12	50	2	Guillo	Guillo					Maintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol>	1	SBSBRKT-01-0000	1	Needs Attention	t	0
19	82	7	Curl #1	Curl #!						2	SBSSD75-04-0000 	6	Needs Attention	t	0
20	88	7	Full Notch	Full Notch						1	SBSCH-05-0000	5	Needs Attention	t	0
18	75	7	Tie Down Slot	Tie Down Slot						1	SBSPL75-04-0000 	4	Needs Attention	t	0
\.


--
-- Name: component_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('component_id_seq', 249, true);


--
-- Data for Name: component_part; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY component_part (component_id, part_id, qty) FROM stdin;
1	1	1
1	2	1
1	3	1
1	4	2
1	5	1
1	6	1
2	7	2
2	8	2
3	9	2
3	10	2
4	11	2
4	12	2
6	13	6
6	14	6
7	15	3
7	16	3
8	17	1
8	18	1
10	19	1
10	20	1
11	21	1
11	22	1
12	23	1
12	24	1
12	25	1
12	26	1
12	27	1
13	28	1
13	29	1
13	30	1
14	31	2
14	32	2
15	33	1
15	34	1
16	35	1
16	36	1
16	37	2
17	38	1
17	39	1
19	40	2
19	41	1
19	42	1
20	43	1
20	44	1
20	45	1
21	46	2
21	47	2
22	48	8
22	49	8
23	50	2
23	51	2
25	52	1
25	53	1
25	54	1
26	55	3
26	56	3
26	57	3
27	58	1
27	59	1
28	60	1
28	61	1
29	62	1
29	63	1
30	64	1
30	65	1
30	66	1
31	67	1
31	68	1
31	69	1
32	70	1
32	71	1
32	72	1
33	73	1
33	74	1
33	75	1
34	76	3
35	77	1
36	78	1
37	79	1
38	80	1
39	81	1
40	82	1
40	83	1
41	84	8
43	85	1
43	86	1
43	87	2
43	88	2
44	89	1
44	90	1
45	91	1
46	92	1
46	93	1
46	94	1
46	95	1
47	96	1
47	97	1
47	98	1
48	99	1
48	100	1
49	101	1
49	102	1
50	103	1
50	104	1
51	105	2
51	106	2
52	107	8
52	108	8
53	109	2
53	110	2
55	111	1
55	112	1
55	113	1
56	114	2
56	115	2
57	116	1
57	117	1
57	118	2
58	119	1
58	120	1
59	121	1
59	122	2
59	123	1
60	55	3
60	56	3
60	57	3
61	58	1
61	59	1
62	60	1
62	61	1
63	62	1
63	63	1
64	64	1
64	65	1
64	66	1
65	67	1
65	68	1
65	69	1
66	70	1
66	71	1
66	72	1
67	73	1
67	74	1
67	75	1
68	82	1
68	83	1
69	84	8
71	85	1
71	86	1
71	87	2
71	88	2
72	31	2
72	32	2
73	33	1
73	34	1
74	35	1
74	36	1
74	37	2
75	38	1
75	39	1
77	40	2
77	41	1
77	42	1
78	43	1
78	44	1
78	45	1
79	46	2
79	47	2
80	48	8
80	49	8
81	50	2
81	51	2
83	52	1
83	53	1
83	54	1
84	55	1
84	56	1
84	57	1
85	58	1
85	59	1
86	60	1
86	61	1
87	62	1
87	63	1
88	64	1
88	65	1
88	66	1
89	67	1
89	68	1
89	69	1
90	70	1
90	71	1
90	72	1
91	73	1
91	74	1
91	75	1
92	76	3
93	77	1
94	78	1
95	79	1
96	80	1
97	81	1
98	82	1
98	83	1
99	84	8
101	85	1
101	86	1
101	87	2
101	88	2
102	89	1
102	90	1
103	91	1
104	31	2
104	32	2
105	33	1
105	34	1
106	35	1
106	36	1
106	37	2
107	38	1
107	39	1
109	40	2
109	41	1
109	42	1
110	43	1
110	44	1
110	45	1
111	46	2
111	47	2
112	48	8
112	49	8
113	50	2
113	51	2
115	52	1
115	53	1
115	54	1
116	55	1
116	56	1
116	57	1
117	58	1
117	59	1
118	60	1
118	61	1
119	62	1
119	63	1
120	64	1
120	65	1
120	66	1
121	67	1
121	68	1
121	69	1
122	70	1
122	71	1
122	72	1
123	73	1
123	74	1
123	75	1
124	76	3
125	77	1
126	78	1
127	79	1
128	80	1
129	81	1
130	82	1
130	83	1
131	84	8
133	85	1
133	86	1
133	87	2
133	88	2
134	89	1
134	90	1
135	91	1
136	105	2
136	106	2
137	107	8
137	108	8
138	109	2
138	110	2
140	111	1
140	112	1
140	113	1
141	114	2
141	115	2
142	116	1
142	117	1
142	118	2
143	119	1
143	120	1
144	121	1
144	122	2
144	123	1
145	31	2
145	32	2
146	33	1
146	34	1
147	35	1
147	36	1
147	37	2
148	38	1
148	39	1
150	40	2
150	41	1
150	42	1
151	43	1
151	44	1
151	45	1
152	46	2
152	47	2
153	48	8
153	49	8
154	50	2
154	51	2
156	52	1
156	53	1
156	54	1
157	55	3
157	56	3
157	57	3
158	58	1
158	59	1
159	60	1
159	61	1
160	62	1
160	63	1
161	64	1
161	65	1
161	66	1
162	67	1
162	68	1
162	69	1
163	70	1
163	71	1
163	72	1
164	73	1
164	74	1
164	75	1
165	76	3
166	77	1
167	78	1
168	79	1
169	80	1
170	81	1
171	82	1
171	83	1
172	84	8
174	85	1
174	86	1
174	87	2
174	88	2
175	89	1
175	90	1
176	91	1
177	1	1
177	2	1
177	3	1
177	4	2
177	5	1
177	6	1
178	7	2
178	8	2
179	9	2
179	10	2
180	11	2
180	12	2
182	13	6
182	14	6
183	15	3
183	16	3
184	17	1
184	18	1
186	19	1
186	20	1
187	21	1
187	22	1
188	23	1
188	24	1
188	25	1
188	26	1
188	27	1
189	28	1
189	29	1
189	30	1
190	124	2
190	125	1
190	126	2
190	127	1
190	128	1
191	129	2
191	130	2
193	131	8
193	132	8
194	133	4
194	134	4
195	135	1
195	136	1
196	137	1
196	138	1
197	139	1
197	140	1
198	141	1
198	142	1
198	143	1
198	144	1
198	145	1
199	146	1
199	147	1
199	148	1
200	31	2
200	32	2
201	33	1
201	34	1
202	35	1
202	36	1
202	37	2
203	38	1
203	39	1
205	40	2
205	41	1
205	42	1
206	43	1
206	44	1
206	45	1
207	46	2
207	47	2
208	48	8
208	49	8
209	50	2
209	51	2
211	52	1
211	53	1
211	54	1
212	55	3
212	56	3
212	57	3
213	58	1
213	59	1
214	60	1
214	61	1
215	62	1
215	63	1
216	64	1
216	65	1
216	66	1
217	67	1
217	68	1
217	69	1
218	70	1
218	71	1
218	72	1
219	73	1
219	74	1
219	75	1
220	76	3
221	77	1
222	78	1
223	79	1
224	80	1
225	81	1
226	82	1
226	83	1
227	84	8
229	85	1
229	86	1
229	87	2
229	88	2
230	89	1
230	90	1
231	91	1
232	92	1
232	93	1
232	94	1
232	95	1
233	96	1
233	97	1
233	98	1
\.


--
-- Data for Name: doc; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY doc (id, name, filename, worker, sitemgr, contractor, type, ref_id, doc_format, notes, filesize, latest_rev, created, user_id, path) FROM stdin;
\.


--
-- Name: doc_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('doc_id_seq', 86, true);


--
-- Data for Name: doc_rev; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY doc_rev (doc_id, id, revdate, descr, filename, user_id, filesize) FROM stdin;
\.


--
-- Name: doc_rev_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('doc_rev_id_seq', 1, false);


--
-- Data for Name: doc_type; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY doc_type (id, name) FROM stdin;
\.


--
-- Data for Name: event; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY event (id, site_id, type, machine_id, tool_id, tool_type, priority, status, startdate, created_by, allocated_by, allocated_to, completed, labour_cost, material_cost, other_cost, notes) FROM stdin;
1	7	Alert	18	74	Notch	1		2016-04-21 02:28:32.696466+09:30	1	0	0	\N	0.00	0.00	0.00	Problem with Notch tool on Plate  machine.
2	7	Alert	19	81	Single Dimple & Rib	1		2016-04-21 02:29:35.735474+09:30	1	0	0	\N	0.00	0.00	0.00	Problem with Single Dimple & Rib tool on Stud machine.
3	2	Alert	5	16	Notch	1		2016-04-22 18:16:25.797981+09:30	1	0	0	\N	0.00	0.00	0.00	Problem with Notch tool on Plate machine.
4	2	Alert	5	19	Nogging	1		2016-04-22 18:23:34.622546+09:30	29	0	0	\N	0.00	0.00	0.00	Problem with Nogging tool on Plate machine.
\.


--
-- Data for Name: event_doc; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY event_doc (event_id, doc_id, doc_rev_id) FROM stdin;
\.


--
-- Name: event_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('event_id_seq', 4, true);


--
-- Data for Name: event_type; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY event_type (id, name) FROM stdin;
\.


--
-- Data for Name: hashtag; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY hashtag (id, name, descr) FROM stdin;
1	SOP33	this is SOP33\n- an item\n- another item\n- one more item
7	SOP22	added some notes
\.


--
-- Name: hashtag_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('hashtag_id_seq', 10, true);


--
-- Data for Name: machine; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY machine (id, site_id, name, descr, make, model, serialnum, is_running, stopped_at, started_at, picture, status, notes, alert_at, electrical, hydraulic, printer, console, rollbed, uncoiler, lube, tasks_to, alerts_to, part_class) FROM stdin;
2	1	Mill	Mill				f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	18
22	8	Plate	Plate 			SBSPL-020	f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	19
37	9	Floor - C250	Floor - C250			SBSFLR-027	t	\N	2015-12-15 08:54:30.116128+10:30		Running		2015-12-14 12:16:58.680011+10:30	Running	Running	Running	Running	Running	Running	Running	0	0	26
1	1	Lathe	Lathe				f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	16
13	4	Wall	Wall			SBSWL-500	f	\N	\N		Needs Attention		2016-03-09 12:59:09.218973+10:30	Running	Running	Running	Running	Running	Running	Running	0	0	11
24	8	Chord	Chord			SBSCH-020	t	2016-01-28 08:56:02.448721+10:30	2016-01-28 10:01:55.939134+10:30		Maintenance Pending		2016-01-29 12:55:06.965232+10:30	Running	Running	Running	Running	Running	Running	Running	0	0	3
3	1	Surface Grinder	Surface Grinder				f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	6
4	2	Floor	Floor				t	2015-12-16 09:52:56.829549+10:30	2016-02-22 15:28:08.661437+10:30		Needs Attention		2016-03-09 12:54:24.584794+10:30	Needs Attention	Running	Running	Running	Running	Stopped	Running	0	0	17
5	2	Plate	Plate			SBSPL-023	t	2015-12-14 13:35:22.337345+10:30	2015-12-14 14:57:23.382371+10:30		Needs Attention		2016-04-22 18:23:34.62733+09:30	Running	Stopped	Running	Running	Running	Running	Running	0	0	19
6	2	Stud	Stud			SBSSD-023	t	2015-12-14 14:44:15.663973+10:30	2015-12-15 08:25:30.045897+10:30		Needs Attention		2016-03-09 12:52:54.24322+10:30	Running	Running	Needs Attention	Running	Stopped	Running	Running	0	0	9
8	2	Web	Web			SBSWB-023	t	\N	2015-12-14 12:01:37.094571+10:30		Running		2016-02-26 12:49:29.113054+10:30	Running	Running	Running	Running	Needs Attention	Running	Running	29	0	8
7	2	Chord	Chord machine desc			SBSCH-023	t	2016-02-08 14:51:07.688473+10:30	2016-02-22 15:28:03.15719+10:30		Running	add some notes to the chord machine	2016-02-23 12:03:34.476895+10:30	Running	Running	Running	Stopped	Running	Running	Needs Attention	0	0	3
9	2	Top Hat 40mm	Top Hat 40mm			SBSTH40-023	t	2015-12-14 11:32:17.763329+10:30	2016-02-22 15:28:15.744571+10:30		Needs Attention		2016-03-09 12:04:49.775348+10:30	Running	Running	Running	Running	Running	Running	Running	0	20	7
10	2	Top Hat 22mm	Top Hat 22mm			SBSTH22-023	f	2016-02-22 15:45:35.184694+10:30	2015-12-15 15:59:15.780496+10:30		Needs Attention		2016-03-09 12:05:41.150922+10:30	Running	Running	Running	Running	Running	Running	Running	0	0	5
11	2	Valley/Lintel	Valley/Lintel			SBSVL-023	t	2015-12-14 14:07:19.773513+10:30	2015-12-14 14:07:34.381414+10:30		Needs Attention		2016-03-09 12:53:41.397503+10:30	Running	Running	Running	Running	Running	Needs Attention	Stopped	0	0	14
12	2	Bracket B4/5	Bracket B4/5			SBSBRKT-023	t	2016-01-28 10:03:07.322981+10:30	2016-02-22 15:27:55.926401+10:30		Needs Attention	Some notes here for the Bracket Machine\nDo this\nDo that\nDo the other\nDo some more	2016-03-24 03:41:23.701566+10:30	Stopped	Needs Attention	Stopped	Needs Attention	Running	Running	Running	0	0	2
14	4	Chord	Chord			SBSCH-500	t	2015-12-23 10:29:25.548214+10:30	2016-01-28 10:55:11.970661+10:30		Needs Attention		2016-01-29 12:19:24.709098+10:30	Running	Running	Running	Running	Running	Running	Running	0	0	3
15	4	Web	Web			SBSWB-500	t	\N	2016-03-07 13:41:41.241799+10:30		Running		\N	Running	Running	Running	Running	Running	Running	Running	0	0	8
16	3	Guillotine	Guillotine				t	\N	2015-12-22 10:19:53.482801+10:30		Running		\N	Running	Running	Running	Running	Running	Running	Running	0	0	15
17	3	Folder	Folder				t	\N	2015-12-22 10:19:45.949302+10:30		Running		\N	Running	Running	Running	Running	Running	Running	Running	0	0	20
18	7	Plate	Plate 			SBSPL-024	t	\N	2015-12-14 15:07:00.863494+10:30		Needs Attention		2016-04-21 02:28:32.727048+09:30	Running	Needs Attention	Running	Running	Running	Running	Running	0	0	19
19	7	Stud	Stud			SBSSD-024	t	\N	2015-12-15 10:24:37.499675+10:30		Needs Attention		2016-04-21 02:29:35.739976+09:30	Running	Running	Running	Running	Running	Running	Running	0	0	9
20	7	Chord	Chord			SBSCH-021	t	2015-12-15 10:42:17.731272+10:30	2015-12-22 10:19:40.332261+10:30		Needs Attention		2016-04-21 01:03:53.779375+09:30	Running	Running	Needs Attention	Running	Running	Running	Running	34	26	3
21	7	Web	Web			SBSWB-021	t	\N	2015-12-14 15:07:12.02879+10:30		Needs Attention		2016-04-21 01:05:04.910987+09:30	Running	Running	Needs Attention	Running	Running	Running	Running	0	0	8
23	8	Stud	Stud			SBSSD-020	f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	9
25	8	Web	Web			SBSWB-020	f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	8
26	8	Wall	Wall			SBSWL-021	f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	11
27	6	Plate	Plate 			SBSPL-028	t	\N	2016-01-12 16:12:22.012583+10:30		Running		\N	Running	Running	Running	Running	Running	Running	Running	0	0	19
28	6	Stud	Stud			SBSSD-028	t	\N	2016-01-12 16:12:28.354697+10:30		Running		\N	Running	Running	Running	Running	Running	Running	Running	0	0	9
30	6	Web	Web			SBSWB-028	t	\N	2016-01-12 16:12:34.195917+10:30		Running		\N	Running	Running	Running	Running	Running	Running	Running	0	0	8
31	10	Lathe	Lathe				f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	16
32	10	Mill	Mill				f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	18
33	10	Surface Grinder	Surface Grinder				f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	6
35	10	Horizontal Bandsaw	Horizontal Bandsaw				f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	13
38	9	Plate	Plate 			SBSPL-027	f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	19
39	9	Stud	Stud			SBSSD-027	f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	9
40	9	Chord	Chord			SBSCH-027	f	2016-01-28 09:27:10.816856+10:30	2015-12-22 10:19:34.514065+10:30		Stopped		\N	Running	Running	Running	Running	Running	Running	Running	0	0	3
41	9	Web	Web			SBSWB-027	f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	8
42	9	Top Hat 40mm	Top Hat 40mm			SBSTH40-027	f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	7
43	9	Top Hat 22mm	Top Hat 22mm			SBSTH22-027	f	\N	\N		New		\N	Running	Running	Running	Running	Running	Running	Running	0	0	5
44	11	Guillotine	Guillotine				t	\N	2015-12-22 10:19:59.684912+10:30		Running		\N	Running	Running	Running	Running	Running	Running	Running	0	0	15
45	11	Folder	Folder				t	2015-12-02 15:51:27.581705+10:30	2015-12-15 09:11:23.620937+10:30		Running		2015-12-14 12:28:16.618377+10:30	Running	Running	Running	Running	Running	Running	Running	0	0	20
36	9	Floor - C150	Floor - C150			SBSFLR-027	t	\N	2015-12-08 13:57:29.146737+10:30		Maintenance Pending		2016-02-02 14:24:58.398784+10:30	Running	Running	Running	Running	Running	Running	Running	0	0	25
29	6	Chord	Chord			SBSCH-028	t	2015-12-23 10:58:43.06844+10:30	2016-01-27 15:05:20.754444+10:30		Needs Attention		2016-03-24 03:39:19.498125+10:30	Running	Running	Running	Running	Running	Running	Running	0	0	3
\.


--
-- Name: machine_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('machine_id_seq', 49, true);


--
-- Data for Name: part; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY part (id, name, descr, stock_code, reorder_stocklevel, reorder_qty, latest_price, qty_type, picture, notes, class, last_price_date, current_stock) FROM stdin;
68	Exit Die Steel	Exit Die Steel - Description	SBSCH-06-1004	1.00	1.00	0.00	ea			3	\N	0.00
69	Punch	Punch - Description	SBSCH-06-1012	1.00	1.00	0.00	ea			3	\N	0.00
1	Notch Die - Rear	Notch Die - Rear - Description	SBSFF150-01-101	1.00	1.00	0.00	ea			17	\N	0.00
2	Notch Die - Front	Notch Die - Front - Description	SBSFF150-01-102	1.00	1.00	0.00	ea			17	\N	0.00
3	Centre Die	Centre Die - Description	SBSFF150-01-103	1.00	1.00	0.00	ea			17	\N	0.00
4	Notch Punch	Notch Punch - Description	SBSFF150-01-113	1.00	1.00	0.00	ea			17	\N	0.00
5	Centre Punch	Centre Punch - Description	SBSFF150-01-114	1.00	1.00	0.00	ea			17	\N	0.00
6	Centre Punch Insert	Centre Punch Insert - Description	SBSFF150-01-115	1.00	1.00	0.00	ea			17	\N	0.00
7	Die Block	Die Block - Description	SBSFF150-02-104	1.00	1.00	0.00	ea			17	\N	0.00
8	Punch	Punch - Description	SBSFF150-02-108	1.00	1.00	0.00	ea			17	\N	0.00
9	Pier Slot Die Block	Pier Slot Die Block - Description	SBSFF150-04-102	1.00	1.00	0.00	ea			17	\N	0.00
76	Entry Guide	Entry Guide - Description	SBSCHTF-02-1005	1.00	1.00	0.00	ea			0	\N	0.00
10	Pier Slot Punch	Pier Slot Punch - Description	SBSFF150-04-106	1.00	1.00	0.00	ea			17	\N	0.00
11	Die Block - Tie Down	Die Block - Tie Down - Description	SBSFF150-23-103	1.00	1.00	0.00	ea			17	\N	0.00
12	Punch - Tie Down	Punch - Tie Down - Description	SBSFF150-23-107	1.00	1.00	0.00	ea			17	\N	0.00
13	Dimple Joist - Punch 	Dimple Joist - Punch  - Description	SBSFF150-08-118	1.00	1.00	0.00	ea			17	\N	0.00
14	Dimple Joist - Die Brush 	Dimple Joist - Die Brush  - Description	SBSFF150-08-120	1.00	1.00	0.00	ea			17	\N	0.00
15	Die Brush	Die Brush - Description	SBSFF150-10-104	1.00	1.00	0.00	ea			17	\N	0.00
16	Punch	Punch - Description	SBSFF150-10-112	1.00	1.00	0.00	ea			17	\N	0.00
17	Service Hole - Die Brush	Service Hole - Die Brush - Description	SBSFF150-12-105	1.00	1.00	0.00	ea			17	\N	0.00
27	Blade Insert - Guillo Joist 	Blade Insert - Guillo Joist  - Description	SBSFF150-19-119	0.00	0.00	0.00	ea			17	\N	0.00
40	Expanding Punch Insert	Expanding Punch Insert - Description	SBSPL75-06-1000	1.00	1.00	0.00	ea			19	\N	0.00
100	Valley Lower Crop Blade	Valley Lower Crop Blade - Description	SBSVL-01-1212	1.00	1.00	0.00	ea			14	\N	0.00
136	Service Hole - Upper Punch	Service Hole - Upper Punch - Description	SBSFF250-12-112	1.00	1.00	0.00	ea			26	\N	0.00
145	Cutting Blade Insert	Cutting Blade Insert - Description	SBSFF250-19-120	1.00	1.00	0.00	ea			26	\N	0.00
90	Exit Guide	Exit Guide - Description	SBSWBTF-01-1007	1.00	1.00	889.00	ea			0	\N	0.00
77	Entry Guide	Entry Guide - Description	SBSCHTF-03-1010	1.00	1.00	0.00	ea			0	\N	0.00
78	Entry Guide	Entry Guide - Description	SBSCHTF-04-1009	1.00	1.00	0.00	ea			0	\N	0.00
79	Entry Guide	Entry Guide - Description	SBSCHTF-05-1006	1.00	1.00	0.00	ea			0	\N	0.00
80	Entry Guide	Entry Guide - Description	SBSCHTF-06-1011	1.00	1.00	0.00	ea			0	\N	0.00
81	Entry Guide	Entry Guide - Description	SBSCHTF-07-1010	1.00	1.00	0.00	ea			0	\N	0.00
82	Turn Over	Turn Over - Description	SBSRB-1171	1.00	1.00	0.00	ea			0	\N	0.00
83	Stitch Wheel	Stitch Wheel - Description	SBSRB-1187	1.00	1.00	0.00	ea			0	\N	0.00
89	Entry Guide	Entry Guide - Description	SBSWBTF-01-1006	1.00	1.00	0.00	ea			0	\N	0.00
91	Entry Guide	Entry Guide - Description	SBSWBTF-02-1009	1.00	1.00	0.00	ea			0	\N	0.00
18	Service Hole - Upper Punch 	Service Hole - Upper Punch  - Description	SBSFF150-12-112	1.00	1.00	0.00	ea			17	\N	0.00
19	Die Brush - Up Dimple	Die Brush - Up Dimple - Description	SBSFF150-17-103	1.00	1.00	0.00	ea			17	\N	0.00
20	Punch - Up Dimple 	Punch - Up Dimple  - Description	SBSFF150-17-111	1.00	1.00	0.00	ea			17	\N	0.00
21	Die Brush - Down Dimple	Die Brush - Down Dimple - Description	SBSFF150-15-103	1.00	1.00	0.00	ea			17	\N	0.00
22	Punch - Down Dimple	Punch - Down Dimple - Description	SBSFF150-15-111	1.00	1.00	0.00	ea			17	\N	0.00
23	Die Steel RH Foot	Die Steel RH Foot - Description	SBSFF150-19-106	1.00	1.00	0.00	ea			17	\N	0.00
24	Die Steel RH Rear 	Die Steel RH Rear  - Description	SBSFF150-19-107	1.00	1.00	0.00	ea			17	\N	0.00
25	Die Steel LH Front	Die Steel LH Front - Description	SBSFF150-19-108	1.00	1.00	0.00	ea			17	\N	0.00
26	Die Steel LH Rear 	Die Steel LH Rear  - Description	SBSFF150-19-109	1.00	1.00	0.00	ea			17	\N	0.00
29	Die Steel LH - Guillo Bearer	Die Steel LH - Guillo Bearer - Description	SBSFF150-20-104	1.00	1.00	0.00	ea			17	\N	0.00
31	Die Bush	Die Bush - Description	SBSPL75-01-1014	1.00	1.00	0.00	ea			19	\N	0.00
32	Punch	Punch - Description	SBSPL75-01-1017	1.00	1.00	0.00	ea			19	\N	0.00
33	Form Punch & Die	Form Punch & Die - Description	SBSPL75-02-1006	1.00	1.00	0.00	ea			19	\N	0.00
34	Punch	Punch - Description	SBSPL75-02-1012	1.00	1.00	0.00	ea			19	\N	0.00
35	Die High Side	Die High Side - Description	SBSPL75-03-1002	1.00	1.00	0.00	ea			19	\N	0.00
36	Die Low Side	Die Low Side - Description	SBSPL75-03-1014	1.00	1.00	0.00	ea			19	\N	0.00
37	Punch	Punch - Description	SBSPL75-03-1013	1.00	1.00	0.00	ea			19	\N	0.00
38	Die Block	Die Block - Description	SBSPL75-04-1003	1.00	1.00	0.00	ea			19	\N	0.00
39	Punch	Punch - Description	SBSPL75-04-1010	1.00	1.00	0.00	ea			19	\N	0.00
41	Die Steel	Die Steel - Description	SBSPL75-06-1002	1.00	1.00	0.00	ea			19	\N	0.00
42	Centre Punch	Centre Punch - Description	SBSPL75-06-1003	1.00	1.00	0.00	ea			19	\N	0.00
43	RH Die Block	RH Die Block - Description	SBSPL75-07-1006	1.00	1.00	0.00	ea			19	\N	0.00
44	Crop Blade	Crop Blade - Description	SBSPL75-07-1012	1.00	1.00	0.00	ea			19	\N	0.00
45	LH Die Block	LH Die Block - Description	SBSPL75-07-1014	1.00	1.00	0.00	ea			19	\N	0.00
46	Die Bush	Die Bush - Description	SBSSD75-01-1006	1.00	1.00	0.00	ea			9	\N	0.00
47	Punch	Punch - Description	SBSSD75-01-1012	1.00	1.00	0.00	ea			9	\N	0.00
48	Punch	Punch - Description	SBSSD75-02-1009	1.00	1.00	0.00	ea			9	\N	0.00
49	Die Bush	Die Bush - Description	SBSSD75-02-1014	1.00	1.00	0.00	ea			9	\N	0.00
50	Punch	Punch - Description	SBSSD75-03-1008	1.00	1.00	0.00	ea			9	\N	0.00
51	Die Bush	Die Bush - Description	SBSSD75-03-1014	1.00	1.00	0.00	ea			9	\N	0.00
52	LH Die Block	LH Die Block - Description	SBSSD75-05-1003	1.00	1.00	0.00	ea			9	\N	0.00
53	RH Die Block	RH Die Block - Description	SBSSD75-05-1004	1.00	1.00	0.00	ea			9	\N	0.00
54	Crop Blade	Crop Blade - Description	SBSSD75-05-1005	1.00	1.00	0.00	ea			9	\N	0.00
55	Die Bush	Die Bush - Description	SBSCH-01-1002	1.00	1.00	0.00	ea			3	\N	0.00
56	Slug Tube	Slug Tube - Description	SBSCH-01-1003	1.00	1.00	0.00	ea			3	\N	0.00
57	Punch	Punch - Description	SBSCH-01-1011	1.00	1.00	0.00	ea			3	\N	0.00
58	Die Block	Die Block - Description	SBSCH-02-1001	1.00	1.00	0.00	ea			3	\N	0.00
59	Punch	Punch - Description	SBSCH-02-1006	1.00	1.00	0.00	ea			3	\N	0.00
60	Die Block	Die Block - Description	SBSCH-03-1003	1.00	1.00	0.00	ea			3	\N	0.00
61	Punch	Punch - Description	SBSCH-03-1009	1.00	1.00	0.00	ea			3	\N	0.00
62	Die Bush	Die Bush - Description	SBSCH-04-1002	1.00	1.00	0.00	ea			3	\N	0.00
63	Punch	Punch - Description	SBSCH-04-1011	1.00	1.00	0.00	ea			3	\N	0.00
64	Die Steel	Die Steel - Description	SBSCH-05-1002	1.00	1.00	0.00	ea			3	\N	0.00
65	Rear Punch	Rear Punch - Description	SBSCH-05-1004	1.00	1.00	0.00	ea			3	\N	0.00
66	Top Punch	Top Punch - Description	SBSCH-05-1014	1.00	1.00	0.00	ea			3	\N	0.00
113	Crop Blade	Crop Blade - Description	SBSSD92-05-1013	1.00	1.00	480.00	ea			11	\N	0.00
138	Punch - Down Dimple	Punch - Down Dimple - Description	SBSFF250-15-111	1.00	1.00	0.00	ea			26	\N	0.00
140	Punch - Up Dimple	Punch - Up Dimple - Description	SBSFF250-17-112	1.00	1.00	0.00	ea			26	\N	0.00
141	Die Steel RH Front	Die Steel RH Front - Description	SBSFF250-19-105	1.00	1.00	0.00	ea			26	\N	0.00
142	Die Steel RH Rear	Die Steel RH Rear - Description	SBSFF250-19-106	1.00	1.00	0.00	ea			26	\N	0.00
143	Die Steel LH Front	Die Steel LH Front - Description	SBSFF250-19-107	1.00	1.00	0.00	ea			26	\N	0.00
144	Die Steel LH Rear	Die Steel LH Rear - Description	SBSFF250-19-108	1.00	1.00	0.00	ea			26	\N	0.00
146	Guillo Die Steel RH	Guillo Die Steel RH - Description	SBSFF250-20-103	1.00	1.00	0.00	ea			26	\N	0.00
28	Die Steel RH - Guillo Bearer	Die Steel RH - Guillo Bearer - Description	SBSFF150-20-103	1.00	1.00	0.00	ea			17	\N	0.00
67	Entry Die Steel	Entry Die Steel - Description	SBSCH-06-1003	1.00	1.00	0.00	ea			3	\N	0.00
70	Entry Die Steel	Entry Die Steel - Description	SBSCH-07-1003	1.00	1.00	0.00	ea			3	\N	0.00
71	Exit Die Steel	Exit Die Steel - Description	SBSCH-07-1004	1.00	1.00	0.00	ea			3	\N	0.00
72	Blade	Blade - Description unchanged	SBSCH-07-1011	23.00	34.00	111.00	per metre		added some notes\nis now part of the bracket machine set	3	2016-04-24	3.10
73	Entry Die Steel	Entry Die Steel - Description	SBSCH-08-1003	1.00	1.00	0.00	ea			3	\N	0.00
74	Exit Die Steel	Exit Die Steel - Description	SBSCH-08-1004	1.00	1.00	0.00	ea			3	\N	0.00
75	Punch	Punch - Description	SBSCH-08-1012	1.00	1.00	0.00	ea			3	\N	0.00
84	Punch (Dayton)	Punch (Dayton) - Description	SBSWB-01-1005	1.00	1.00	0.00	ea			8	\N	0.00
85	Form & Pierce Die	Form & Pierce Die - Description	SBSWB-03-1011	1.00	1.00	397.00	ea			8	\N	0.00
86	Punch - Bow Tie	Punch - Bow Tie - Description	SBSWB-03-1013	1.00	1.00	499.00	ea			8	\N	0.00
87	Die Bush Matrix	Die Bush Matrix - Description	SBSWB-03-1016	1.00	1.00	328.00	ea			8	\N	0.00
88	Punch - Dayton	Punch - Dayton - Description	SBSWB-03-1017	1.00	1.00	0.00	ea			8	\N	0.00
92	LH Die Block	LH Die Block - Description	SBSPM-01-103	1.00	1.00	0.00	ea			7	\N	0.00
93	RH Die Block	RH Die Block - Description	SBSPM-01-104	1.00	1.00	0.00	ea			7	\N	0.00
94	Punch	Punch - Description	SBSPM-01-112	1.00	1.00	0.00	ea			7	\N	0.00
95	Punch Insert	Punch Insert - Description	SBSPM-01-113	1.00	1.00	0.00	ea			7	\N	0.00
96	RH Die Block	RH Die Block - Description	SBSTH22-01-1007	1.00	1.00	0.00	ea			5	\N	0.00
97	LH Die Block	LH Die Block - Description	SBSTH22-01-1008	1.00	1.00	0.00	ea			5	\N	0.00
98	Crop Blade	Crop Blade - Description	SBSTH22-01-1009	1.00	1.00	0.00	ea			5	\N	0.00
99	Valley Upper Crop Balde	Valley Upper Crop Balde - Description	SBSVL-01-1208	1.00	1.00	0.00	ea			14	\N	0.00
101	Lintel Lower Crop Blade	Lintel Lower Crop Blade - Description	SBSVL-01-1236	1.00	1.00	0.00	ea			14	\N	0.00
102	Lintel Upper Crop Blade	Lintel Upper Crop Blade - Description	SBSVL-01-1240	1.00	1.00	0.00	ea			14	\N	0.00
103	Upper Crop Blade	Upper Crop Blade - Description	SBSBRKT-01-1005	1.00	1.00	0.00	ea			2	\N	0.00
104	Lower Crop Blade	Lower Crop Blade - Description	SBSBRKT-01-1006	1.00	1.00	0.00	ea			2	\N	0.00
105	Form Punch & Die	Form Punch & Die - Description	SBSSD92-01-1004	1.00	1.00	0.00	ea			11	\N	0.00
106	Slot Punch	Slot Punch - Description	SBSSD92-01-1013	1.00	1.00	0.00	ea			11	\N	0.00
107	Punch	Punch - Description	SBSSD92-02-1009	1.00	1.00	0.00	ea			11	\N	0.00
108	Die Bush	Die Bush - Description	SBSSD92-02-1014	1.00	1.00	0.00	ea			11	\N	0.00
109	Punch	Punch - Description	SBSSD92-03-1011	1.00	1.00	0.00	ea			11	\N	0.00
110	Die Bush	Die Bush - Description	SBSSD92-03-1015	1.00	1.00	0.00	ea			11	\N	0.00
111	Die Steel RH	Die Steel RH - Description	SBSSD92-05-1002	1.00	1.00	0.00	ea			11	\N	0.00
112	Die Steel LH	Die Steel LH - Description	SBSSD92-05-1003	1.00	1.00	0.00	ea			11	\N	0.00
114	Punch	Punch - Description	SBSPL92-01-1011	1.00	1.00	0.00	ea			11	\N	0.00
115	Die Bush	Die Bush - Description	SBSPL92-01-1015	1.00	1.00	0.00	ea			11	\N	0.00
116	Die High Side	Die High Side - Description	SBSPL92-03-1004	1.00	1.00	0.00	ea			11	\N	0.00
117	Die Low Side	Die Low Side - Description	SBSPL92-03-1005	1.00	1.00	0.00	ea			11	\N	0.00
118	Punch	Punch - Description	SBSPL92-03-1009	1.00	1.00	0.00	ea			11	\N	0.00
119	Die Block	Die Block - Description	SBSPL92-04-1005	1.00	1.00	0.00	ea			11	\N	0.00
120	Punch	Punch - Description	SBSPL92-04-1010	1.00	1.00	0.00	ea			11	\N	0.00
121	Die Steel	Die Steel - Description	SBSPL92-06-1003	1.00	1.00	0.00	ea			11	\N	0.00
122	Punch	Punch - Description	SBSPL92-06-1012	1.00	1.00	0.00	ea			11	\N	0.00
123	Punch Expander	Punch Expander - Description	SBSPL92-06-1014	1.00	1.00	0.00	ea			11	\N	0.00
124	Notch Die - H-Cut	Notch Die - H-Cut - Description	SBSFF250-01-101	1.00	1.00	0.00	ea			26	\N	0.00
125	Centre Die - H-Cut	Centre Die - H-Cut - Description	SBSFF250-01-102	1.00	1.00	0.00	ea			26	\N	0.00
126	Notch Punch - H-Cut	Notch Punch - H-Cut - Description	SBSFF250-01-111	1.00	1.00	0.00	ea			26	\N	0.00
127	Centre Punch - H-Cut	Centre Punch - H-Cut - Description	SBSFF250-01-112	1.00	1.00	0.00	ea			26	\N	0.00
128	Centre Punch Insert - H-Cut	Centre Punch Insert - H-Cut - Description	SBSFF250-01-113	1.00	1.00	111.00	ea			26	\N	33.00
129	Notch Die - Notch 	Notch Die - Notch  - Description	SBSFF250-02-101	1.00	1.00	0.00	ea			26	\N	0.00
130	Notch Punch - Notch 	Notch Punch - Notch  - Description	SBSFF250-02-108	1.00	1.00	0.00	ea			26	\N	0.00
131	Punch - Dimple Joist 	Punch - Dimple Joist  - Description	SBSFF250-08-115	1.00	1.00	0.00	ea			26	\N	0.00
132	Die Bush - Dimple Joist	Die Bush - Dimple Joist - Description	SBSFF250-08-117	1.00	1.00	0.00	ea			26	\N	0.00
133	Dimple Die - Dimple Bearer	Dimple Die - Dimple Bearer - Description	SBSFF250-10-103	1.00	1.00	0.00	ea			26	\N	0.00
134	Punch - Dimple Bearer	Punch - Dimple Bearer - Description	SBSFF250-10-111	1.00	1.00	0.00	ea			26	\N	0.00
135	Service Hole - Die Brush	Service Hole - Die Brush - Description	SBSFF250-12-105	1.00	1.00	0.00	ea			26	\N	0.00
137	Bottom Die - Down Dimple	Bottom Die - Down Dimple - Description	SBSFF250-15-103	1.00	1.00	2.00	ea			26	2016-04-24	1.00
30	Blade - Guillo Bearer	Blade - Guillo Bearer - Description	SBSFF150-20-108	0.00	0.00	0.00	ea			17	2016-04-25	0.00
139	Bottom Die - Up Dimple	Bottom Die - Up Dimple - Description	SBSFF250-17-104	1.00	1.00	33.00	ea			26	2016-04-24	222.00
147	Guillo Die Steel LH	Guillo Die Steel LH - Description	SBSFF250-20-104	1.00	1.00	0.00	ea			26	\N	0.00
148	Guillo Straight Blade	Guillo Straight Blade - Description	SBSFF250-20-108	1.00	1.00	0.00	ea			26	\N	0.00
\.


--
-- Data for Name: part_class; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY part_class (id, name, descr) FROM stdin;
2	Bracket B4/5	
3	Chord	
5	Top Hat 22mm	
6	Surface Grinder	
7	Top Hat 40mm	
8	Web	
9	Stud	
11	Wall	
13	Horizontal Bandsaw	
14	Valley/Lintel	
15	Guillotine	
16	Lathe	
18	Mill	
19	Plate	
20	Folder	
26	Floor - C250	Floor 250 Connecticut
17	Floor	Floor Machine - includes C150
\.


--
-- Name: part_class_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('part_class_id_seq', 26, true);


--
-- Name: part_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('part_id_seq', 159, true);


--
-- Data for Name: part_price; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY part_price (part_id, datefrom, price) FROM stdin;
128	2016-04-24 22:03:14.789412+09:30	111.00
139	2016-04-24 22:04:57.373549+09:30	33.00
72	2016-04-24 23:13:25.343753+09:30	111.00
137	2016-04-24 23:16:09.746906+09:30	2.00
30	2016-04-25 05:09:39.25151+09:30	1.00
30	2016-04-25 05:09:51.464428+09:30	0.00
\.


--
-- Data for Name: part_stock; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY part_stock (part_id, datefrom, stock_level) FROM stdin;
159	2016-04-24 20:22:42.74564+09:30	33.12
159	2016-04-24 21:54:19.979042+09:30	45.00
157	2016-04-24 21:57:29.464402+09:30	23.00
157	2016-04-24 21:58:34.349229+09:30	24.15
157	2016-04-24 22:00:49.376168+09:30	24.11
157	2016-04-24 22:00:54.458137+09:30	23.33
157	2016-04-24 22:00:57.07475+09:30	23.12
157	2016-04-24 22:01:31.279983+09:30	21.33
157	2016-04-24 22:02:17.17796+09:30	19.22
128	2016-04-24 22:03:14.785183+09:30	33.00
139	2016-04-24 22:05:34.648088+09:30	222.00
159	2016-04-24 23:05:52.552402+09:30	12.00
72	2016-04-24 23:13:21.295262+09:30	3.00
72	2016-04-24 23:13:30.781796+09:30	3.10
159	2016-04-24 23:14:54.502332+09:30	13.00
137	2016-04-24 23:16:09.738491+09:30	1.00
30	2016-04-25 05:09:39.242923+09:30	1.00
30	2016-04-25 05:09:51.455894+09:30	0.00
\.


--
-- Data for Name: part_vendor; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY part_vendor (part_id, vendor_id, vendor_code, latest_price) FROM stdin;
86	2	BT23	499.00
113	2	LeCutter	480.00
85	2	PD6	397.00
87	2	DBrush7	328.00
\.


--
-- Data for Name: sched_control; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sched_control (id, last_run) FROM stdin;
\.


--
-- Data for Name: sched_control_task; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sched_control_task (task_id, last_gen, last_jobcount) FROM stdin;
\.


--
-- Data for Name: sched_task; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sched_task (id, machine_id, comp_type, tool_id, component, descr, startdate, oneoffdate, freq, parent_task, days, count, week, duration_days, labour_cost, material_cost, other_cost_desc, other_cost, last_generated, weekday, user_id, paused) FROM stdin;
17	8	T	41	Pierce Location 	this should generate immediately	0001-01-01	0001-01-01	Monthly	\N	0	0	3	444	121.00	333.00	\N	\N	2016-04-18	1	0	f
22	6	A	0	General Maintenance	test case friday	\N	\N	Monthly	\N	\N	\N	3	1	0.00	0.00	\N	\N	2016-04-22	5	0	f
24	6	A	0	General Maintenance	test case monday	\N	\N	Monthly	\N	\N	\N	3	1	0.00	0.00	\N	\N	2016-04-18	1	28	f
16	7	A	0	General Maintenance	this one wont generate yet\nbut will generate now	2016-07-22	\N	Yearly	\N	\N	\N	\N	1	0.00	0.00	\N	\N	\N	\N	0	f
11	7	C	0	Console	every 16 jobs, need to reboot the windows\nand that will take 2 days	\N	\N	Job Count	\N	\N	16	\N	2	10.00	7.00	\N	\N	\N	\N	0	f
15	7	A	0	General Maintenance	second week maint	\N	\N	Monthly	\N	\N	\N	2	1	0.00	0.00	\N	\N	\N	1	0	f
23	6	A	0	General Maintenance	test case thursday	0001-01-01	0001-01-01	Monthly	\N	0	0	3	1	0.00	0.00	\N	\N	2016-04-21	4	0	t
18	8	C	0	RollBed	rollbed fix every 3 days	\N	\N	Every N Days	\N	3	\N	\N	1	0.00	0.00	\N	\N	2016-05-10	\N	0	f
19	6	C	0	Lube	relube the stud machine	\N	\N	Every N Days	\N	2	\N	\N	1	0.00	0.00	\N	\N	2016-05-10	\N	0	f
20	16	C	0	Printer	add paper to printer	0001-01-01	0001-01-01	Every N Days	\N	3	0	0	1	0.00	2.00	\N	\N	2016-05-10	\N	0	f
26	29	A	0	General Maintenance	This is a one off task with some checklist items\n- the first check item\n- second check item\n- third item\n- fourth item\n- fifth item	\N	2016-04-25	One Off	\N	\N	\N	\N	1	0.00	0.00	\N	\N	2016-04-25	1	0	f
27	29	A	0	General Maintenance	Every couple of days, do the following actions\n\n- Antarctica\n- America\n- Andalusia\n- Belize	0001-01-01	0001-01-01	Every N Days	\N	2	0	0	1	0.00	0.00	\N	\N	\N	1	0	t
28	29	A	0	General Maintenance	This one goes off every day with the following weapons\n- Machine gun\n- Rifle\n- Chainsaw\n- Atomic Bomb	\N	\N	Every N Days	\N	1	\N	\N	1	0.00	0.00	\N	\N	2016-05-10	1	0	f
29	28	A	0	General Maintenance	A new farm with these animals\n\n- horse\n- cow\n- duck\n- turtle	\N	\N	Every N Days	\N	1	\N	\N	1	0.00	0.00	\N	\N	2016-05-10	1	0	f
13	7	A	0	General Maintenance	last week maint	\N	\N	Monthly	\N	\N	\N	4	1	0.00	0.00	\N	\N	2016-04-25	1	0	f
25	12	A	0	General Maintenance	do something every 4 days	0001-01-01	0001-01-01	Every N Days	\N	4	0	0	1	0.00	0.00	\N	\N	2016-05-07	1	45	f
10	7	T	32	Straight Guillo	straighten the guillo every 44 jobs that are run	\N	\N	Job Count	\N	\N	44	\N	1	0.00	0.00	\N	\N	\N	1	0	f
6	7	T	29	Up Dimple	up dimple on 3rd week of each month	\N	\N	Monthly	\N	\N	\N	3	1	0.00	0.00	\N	\N	2016-04-21	4	0	f
7	7	C	0	Hydraulic	yearly hydraulics overhaul	2016-04-20	\N	Yearly	\N	\N	\N	\N	1	222.00	333.00	\N	\N	2016-04-20	1	0	f
8	7	A	0	General Maintenance	One off general maintenance\n	\N	2016-04-20	One Off	\N	\N	\N	\N	1	0.00	0.00	\N	\N	2016-04-20	1	0	f
9	7	T	31	Right Angle Guillo	every 22 days, fix the RA guillo	\N	\N	Every N Days	\N	22	\N	\N	1	0.00	0.00	\N	\N	2016-04-22	1	0	f
14	7	A	0	General Maintenance	first week maint	\N	\N	Monthly	\N	\N	\N	1	1	0.00	0.00	\N	\N	2016-05-02	1	0	f
21	12	A	0	General Maintenance	Tue maintenance	\N	\N	Monthly	\N	\N	\N	1	1	0.00	0.00	\N	\N	2016-05-02	1	49	f
\.


--
-- Name: sched_task_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('sched_task_id_seq', 29, true);


--
-- Data for Name: sched_task_part; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sched_task_part (task_id, part_id, qty, notes) FROM stdin;
6	58	23.00	first one
7	62	435.00	uouoeu
6	60	444.00	tnoe utnhoeunheou\nmodified
8	58	324.00	oeuoe
7	70	1.00	
7	63	2.00	
10	64	2.00	
9	59	3.00	oeuu
9	71	4.00	oeuoeu
26	72	2.00	2 blade items
26	58	3.00	3 block itiems
26	62	3.00	some die bushes as well
28	68	22.00	calibre
28	64	66.00	die steel
\.


--
-- Data for Name: site; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY site (id, name, address, phone, fax, image, parent_site, notes, stock_site, x, y, alerts_to, tasks_to, manager) FROM stdin;
11	Connecticut Workshop					9		0	0	0	0	0	0
13	Connecticut Thermaloc					9		0	0	0	0	0	0
3	Edinburgh - Fab Shop	40 Barfield Crescent, Edinburgh North, SA 5113	+61 8 8282 7272	+61 8 8282 7251		2		2	0	0	0	0	0
4	Edinburgh - SMiC	40 Barfield Crescent, Edinburgh North, SA 5113	+61 8 8282 7272	+61 8 8282 7251		2		2	0	0	0	0	0
1	Machine Workshop	40 Barfield Crescent, Edinburgh North, SA 5113	+61 8 8282 7272	+61 8 8282 7251		2		2	0	0	0	0	0
7	Minto - Factory	43 - 47 Lincoln Street, Minto, NSW 2566	+61 2 4913 7858	+61 2 4913 7899		8		2	0	0	0	0	0
5	Thermaloc	40 Barfield Crescent, Edinburgh North, SA 5113	+61 8 8282 7272	+61 8 8282 7251		2		2	0	0	0	0	0
6	Tomago - Factory	606 Tomago Road, Tomago, NSW 2322	+61 2 4913 7888	+61 2 4913 7899		0	Tomago Site Level Maintenance Procedure\n\nAll service contactors arriving at Tomago must register with the site safety mgr before commencing work, thanks	2	0	0	0	0	0
9	Connecticut - Factory	Connecticut	TBA	TBA		0		0	0	0	0	0	0
12	Connecticut Fab					9		0	0	0	43	0	0
2	Edinburgh - Factory	40 Barfield Crescent, Edinburgh North, SA 5113	+61 8 8282 7272	+61 8 8282 7251		0	Main Edinburgh Site\nThis includes all machines at the main factory, and any machines at sub-factories at the main Edinburgh site.	0	0	0	8	30	0
8	Chinderah - Factory	14 Ozone Street, Chinderah, NSW 2487 another edit	+61 2 4913 7867	+61 2 6674 5027		0		2	0	0	39	0	39
\.


--
-- Name: site_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('site_id_seq', 16, true);


--
-- Data for Name: site_layout; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY site_layout (site_id, seq, machine_id, span) FROM stdin;
1	1	1	12
1	2	2	12
1	3	3	12
2	3	6	6
2	4	7	6
2	7	11	4
2	8	10	4
2	9	9	4
4	1	13	12
4	2	14	12
4	3	15	12
6	1	27	12
6	2	28	12
6	3	29	12
6	4	30	12
7	1	20	6
7	2	19	6
7	3	21	6
7	4	18	6
2	1	12	4
2	2	5	8
2	5	8	4
2	6	4	8
\.


--
-- Data for Name: skill; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY skill (id, name, notes) FROM stdin;
1	Manager	
2	Supervisor	
3	Technician	
5	Site Manager	
6	Hydraulics	
4	Automation	Test of notes with a link to a course on BSc in Automation &nbsp;dsf<br><br><a href="http://www.uab.cat/web/studying/ehea-degrees/study-plan/skills/x-1345467897139.html?param1=1232089764776">http://www.uab.cat/web/studying/ehea-degrees/study-plan/skills/x-1345467897139.html?param1=1232089764776</a><br>
\.


--
-- Name: skill_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('skill_id_seq', 7, true);


--
-- Data for Name: sm_component; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sm_component (task_id, machine_id, component, completed, mins_spent, labour_cost, materal_cost, notes) FROM stdin;
\.


--
-- Data for Name: sm_component_item; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sm_component_item (task_id, component, seq, notes, done) FROM stdin;
\.


--
-- Data for Name: sm_machine; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sm_machine (task_id, machine_id, completed) FROM stdin;
\.


--
-- Data for Name: sm_machine_item; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sm_machine_item (task_id, machine_id, seq, notes, done) FROM stdin;
\.


--
-- Data for Name: sm_parts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sm_parts (task_id, part_id, date, qty, value) FROM stdin;
\.


--
-- Data for Name: sm_task; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sm_task (id, user_id, date, type, week, completed, escalate_date, escalate_user, status) FROM stdin;
\.


--
-- Name: sm_task_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('sm_task_id_seq', 1, false);


--
-- Data for Name: sm_tool; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sm_tool (task_id, machine_id, tool_id, completed, mins_spent, labour_cost, materal_cost, notes) FROM stdin;
\.


--
-- Data for Name: sm_tool_item; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sm_tool_item (task_id, tool_id, seq, notes, done) FROM stdin;
\.


--
-- Data for Name: stock_level; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY stock_level (part_id, site_id, datefrom, qty, notes) FROM stdin;
\.


--
-- Name: stock_level_part_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('stock_level_part_id_seq', 1, false);


--
-- Data for Name: sys_log; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY sys_log (id, status, type, ref_type, ref_id, logdate, ip, descr, user_id, username, before, after) FROM stdin;
1	1	InitData	I	1	2015-11-28 14:19:07.942218+10:30	localhost	Initialize Data	1	Admin		
\.


--
-- Name: sys_log_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('sys_log_id_seq', 2615, true);


--
-- Data for Name: task; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY task (id, sched_id, machine_id, tool_id, comp_type, component, descr, log, created_date, startdate, due_date, escalate_date, assigned_by, assigned_to, assigned_date, completed_date, has_issue, issue_resolved_date, labour_est, material_est, labour_cost, material_cost, other_cost_desc, other_cost) FROM stdin;
114	6	7	29	T	Up Dimple	up dimple on 3rd week of each month		2016-04-22 13:11:36.450971+09:30	2016-04-21	2016-04-22	2016-05-21	\N	30	2016-04-21 00:00:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
115	7	7	0	C	Hydraulic	yearly hydraulics overhaul		2016-04-22 13:11:36.460533+09:30	2016-04-20	2016-04-21	2016-05-20	\N	30	2016-04-20 09:30:00+09:30	\N	f	\N	222.00	333.00	0.00	0.00	\N	\N
116	8	7	0	A	General Maintenance	One off general maintenance\n		2016-04-22 13:11:36.469814+09:30	2016-04-20	2016-04-21	2016-05-20	\N	30	2016-04-20 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
117	9	7	31	T	Right Angle Guillo	every 22 days, fix the RA guillo		2016-04-22 13:11:36.478874+09:30	2016-04-22	2016-04-23	2016-05-22	\N	30	2016-04-22 13:11:36.444685+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
118	17	8	41	T	Pierce Location 	this should generate immediately		2016-04-22 13:11:36.488461+09:30	2016-04-18	2016-04-22	2016-05-18	\N	29	2016-04-18 00:00:00+09:30	\N	f	\N	121.00	333.00	0.00	0.00	\N	\N
119	18	8	0	C	RollBed	rollbed fix every 3 days		2016-04-22 13:11:36.497647+09:30	2016-04-22	2016-04-23	2016-05-22	\N	29	2016-04-22 13:11:36.444685+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
120	18	8	0	C	RollBed	rollbed fix every 3 days		2016-04-22 13:11:36.507091+09:30	2016-04-25	2016-04-26	2016-05-25	\N	29	2016-04-25 13:11:36.444685+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
121	18	8	0	C	RollBed	rollbed fix every 3 days		2016-04-22 13:11:36.51623+09:30	2016-04-28	2016-04-29	2016-05-28	\N	29	2016-04-28 13:11:36.444685+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
122	19	6	0	C	Lube	relube the stud machine		2016-04-22 13:11:36.525604+09:30	2016-04-22	2016-04-23	2016-05-22	\N	30	2016-04-22 13:11:36.444685+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
123	19	6	0	C	Lube	relube the stud machine		2016-04-22 13:11:36.534901+09:30	2016-04-24	2016-04-25	2016-05-24	\N	30	2016-04-24 13:11:36.444685+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
124	19	6	0	C	Lube	relube the stud machine		2016-04-22 13:11:36.544064+09:30	2016-04-26	2016-04-27	2016-05-26	\N	30	2016-04-26 13:11:36.444685+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
125	19	6	0	C	Lube	relube the stud machine		2016-04-22 13:11:36.553505+09:30	2016-04-28	2016-04-29	2016-05-28	\N	30	2016-04-28 13:11:36.444685+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
126	20	16	0	C	Printer	add paper to printer		2016-04-22 13:11:36.563964+09:30	2016-04-22	2016-04-23	2016-05-22	\N	0	2016-04-22 13:11:36.444685+09:30	\N	f	\N	0.00	2.00	0.00	0.00	\N	\N
127	20	16	0	C	Printer	add paper to printer		2016-04-22 13:11:36.573549+09:30	2016-04-25	2016-04-26	2016-05-25	\N	0	2016-04-25 13:11:36.444685+09:30	\N	f	\N	0.00	2.00	0.00	0.00	\N	\N
128	20	16	0	C	Printer	add paper to printer		2016-04-22 13:11:36.582724+09:30	2016-04-28	2016-04-29	2016-05-28	\N	0	2016-04-28 13:11:36.444685+09:30	\N	f	\N	0.00	2.00	0.00	0.00	\N	\N
129	22	6	0	A	General Maintenance	test case friday		2016-04-22 13:11:36.592188+09:30	2016-04-22	2016-04-22	2016-05-22	\N	30	2016-04-22 00:00:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
130	23	6	0	A	General Maintenance	test case thursday		2016-04-22 13:11:36.601544+09:30	2016-04-21	2016-04-22	2016-05-21	\N	30	2016-04-21 00:00:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
131	24	6	0	A	General Maintenance	test case monday		2016-04-22 13:11:36.61099+09:30	2016-04-18	2016-04-22	2016-05-18	\N	30	2016-04-18 00:00:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
132	6	7	29	T	Up Dimple	up dimple on 3rd week of each month		2016-04-22 16:11:27.228422+09:30	2016-04-21	2016-04-22	2016-05-21	\N	30	2016-04-21 00:00:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
133	6	7	29	T	Up Dimple	up dimple on 3rd week of each month		2016-04-22 18:13:35.642257+09:30	2016-04-21	2016-04-22	2016-05-21	\N	30	2016-04-21 00:00:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
134	19	6	0	C	Lube	relube the stud machine		2016-04-23 09:43:52.249491+09:30	2016-04-30	2016-05-01	2016-05-30	\N	30	2016-04-30 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
135	13	7	0	A	General Maintenance	last week maint		2016-04-24 00:59:01.421686+09:30	2016-04-25	2016-04-29	2016-05-25	\N	30	2016-04-25 00:00:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
136	18	8	0	C	RollBed	rollbed fix every 3 days		2016-04-24 09:59:01.503947+09:30	2016-05-01	2016-05-02	2016-06-01	\N	29	2016-05-01 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
137	20	16	0	C	Printer	add paper to printer		2016-04-24 09:59:01.537048+09:30	2016-05-01	2016-05-02	2016-06-01	\N	0	2016-05-01 09:30:00+09:30	\N	f	\N	0.00	2.00	0.00	0.00	\N	\N
138	19	6	0	C	Lube	relube the stud machine		2016-04-25 09:32:28.536727+09:30	2016-05-02	2016-05-03	2016-06-02	\N	30	2016-05-02 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
139	25	12	0	A	General Maintenance	do something every 4 days		2016-04-25 12:19:58.627859+09:30	2016-04-25	2016-04-26	2016-05-25	\N	30	2016-04-25 12:19:58.626153+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
140	25	12	0	A	General Maintenance	do something every 4 days		2016-04-25 12:19:58.637203+09:30	2016-04-29	2016-04-30	2016-05-29	\N	30	2016-04-29 12:19:58.626153+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
141	8	7	0	A	General Maintenance	One off general maintenance\n		2016-04-25 15:55:32.385842+09:30	2016-04-20	2016-04-21	2016-05-20	\N	30	2016-04-20 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
143	26	29	0	A	General Maintenance	This is a one off task with some checklist items\n- the first check item\n- second check item\n- third item\n- fourth item\n- fifth item		2016-04-25 21:34:25.998782+09:30	2016-04-25	2016-04-26	2016-05-25	\N	0	2016-04-25 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
144	26	29	0	A	General Maintenance	This is a one off task with some checklist items\n- the first check item\n- second check item\n- third item\n- fourth item\n- fifth item		2016-04-25 21:55:24.830555+09:30	2016-04-25	2016-04-26	2016-05-25	\N	0	2016-04-25 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
145	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n- Machine gun\n- Rifle\n- Chainsaw\n- Atomic Bomb		2016-04-25 22:01:15.142717+09:30	2016-04-26	2016-04-27	2016-05-26	\N	0	2016-04-26 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
146	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n- Machine gun\n- Rifle\n- Chainsaw\n- Atomic Bomb		2016-04-25 22:01:15.17799+09:30	2016-04-27	2016-04-28	2016-05-27	\N	0	2016-04-27 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
147	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n- Machine gun\n- Rifle\n- Chainsaw\n- Atomic Bomb		2016-04-25 22:01:15.211865+09:30	2016-04-28	2016-04-29	2016-05-28	\N	0	2016-04-28 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
148	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n- Machine gun\n- Rifle\n- Chainsaw\n- Atomic Bomb		2016-04-25 22:01:15.245685+09:30	2016-04-29	2016-04-30	2016-05-29	\N	0	2016-04-29 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
149	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n- Machine gun\n- Rifle\n- Chainsaw\n- Atomic Bomb		2016-04-25 22:01:15.279516+09:30	2016-04-30	2016-05-01	2016-05-30	\N	0	2016-04-30 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
150	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n- Machine gun\n- Rifle\n- Chainsaw\n- Atomic Bomb		2016-04-25 22:01:15.312745+09:30	2016-05-01	2016-05-02	2016-06-01	\N	0	2016-05-01 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
152	29	28	0	A	General Maintenance	A new farm with these animals\n\n- horse\n- cow\n- duck\n- turtle		2016-04-25 22:05:02.15771+09:30	2016-04-26	2016-04-27	2016-05-26	\N	0	2016-04-26 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
153	29	28	0	A	General Maintenance	A new farm with these animals\n\n- horse\n- cow\n- duck\n- turtle		2016-04-25 22:05:02.184214+09:30	2016-04-27	2016-04-28	2016-05-27	\N	0	2016-04-27 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
154	29	28	0	A	General Maintenance	A new farm with these animals\n\n- horse\n- cow\n- duck\n- turtle		2016-04-25 22:05:02.210395+09:30	2016-04-28	2016-04-29	2016-05-28	\N	0	2016-04-28 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
155	29	28	0	A	General Maintenance	A new farm with these animals\n\n- horse\n- cow\n- duck\n- turtle		2016-04-25 22:05:02.236385+09:30	2016-04-29	2016-04-30	2016-05-29	\N	0	2016-04-29 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
156	29	28	0	A	General Maintenance	A new farm with these animals\n\n- horse\n- cow\n- duck\n- turtle		2016-04-25 22:05:02.261984+09:30	2016-04-30	2016-05-01	2016-05-30	\N	0	2016-04-30 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
157	29	28	0	A	General Maintenance	A new farm with these animals\n\n- horse\n- cow\n- duck\n- turtle		2016-04-25 22:05:02.28726+09:30	2016-05-01	2016-05-02	2016-06-01	\N	0	2016-05-01 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
159	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-25 22:06:44.105373+09:30	2016-04-26	2016-04-27	2016-05-26	\N	0	2016-04-26 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
160	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-25 22:06:44.136392+09:30	2016-04-27	2016-04-28	2016-05-27	\N	0	2016-04-27 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
151	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n- Machine gun\n- Rifle\n- Chainsaw\n- Atomic Bomb		2016-04-25 22:01:15.346266+09:30	2016-05-02	2016-05-03	2016-06-02	\N	0	2016-05-02 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
158	29	28	0	A	General Maintenance	A new farm with these animals\n\n- horse\n- cow\n- duck\n- turtle		2016-04-25 22:05:02.313355+09:30	2016-05-02	2016-05-03	2016-06-02	\N	0	2016-05-02 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
161	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-25 22:06:44.167092+09:30	2016-04-28	2016-04-29	2016-05-28	\N	0	2016-04-28 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
162	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-25 22:06:44.197407+09:30	2016-04-29	2016-04-30	2016-05-29	\N	0	2016-04-29 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
163	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-25 22:06:44.227345+09:30	2016-04-30	2016-05-01	2016-05-30	\N	0	2016-04-30 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
164	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-25 22:06:44.257543+09:30	2016-05-01	2016-05-02	2016-06-01	\N	0	2016-05-01 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
165	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-25 22:06:44.287866+09:30	2016-05-02	2016-05-03	2016-06-02	\N	0	2016-05-02 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
166	25	12	0	A	General Maintenance	do something every 4 days		2016-04-26 10:12:04.623696+09:30	2016-05-03	2016-05-04	2016-06-03	\N	30	2016-05-03 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
167	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n		2016-04-26 10:12:04.633322+09:30	2016-05-03	2016-05-04	2016-06-03	\N	0	2016-05-03 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
168	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-26 10:12:04.671496+09:30	2016-05-03	2016-05-04	2016-06-03	\N	0	2016-05-03 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
169	18	8	0	C	RollBed	rollbed fix every 3 days		2016-04-27 09:52:31.785076+09:30	2016-05-04	2016-05-05	2016-06-04	\N	29	2016-05-04 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
170	19	6	0	C	Lube	relube the stud machine		2016-04-27 09:52:31.818697+09:30	2016-05-04	2016-05-05	2016-06-04	\N	30	2016-05-04 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
171	20	16	0	C	Printer	add paper to printer		2016-04-27 09:52:31.827426+09:30	2016-05-04	2016-05-05	2016-06-04	\N	0	2016-05-04 09:30:00+09:30	\N	f	\N	0.00	2.00	0.00	0.00	\N	\N
172	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n		2016-04-27 09:52:31.835944+09:30	2016-05-04	2016-05-05	2016-06-04	\N	0	2016-05-04 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
173	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-27 09:52:31.878699+09:30	2016-05-04	2016-05-05	2016-06-04	\N	0	2016-05-04 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
174	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n		2016-04-28 09:52:33.379042+09:30	2016-05-05	2016-05-06	2016-06-05	\N	0	2016-05-05 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
175	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-28 09:52:33.422971+09:30	2016-05-05	2016-05-06	2016-06-05	\N	0	2016-05-05 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
176	19	6	0	C	Lube	relube the stud machine		2016-04-29 09:57:59.781239+09:30	2016-05-06	2016-05-07	2016-06-06	\N	30	2016-05-06 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
177	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n		2016-04-29 09:57:59.791378+09:30	2016-05-06	2016-05-07	2016-06-06	\N	0	2016-05-06 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
178	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-29 09:57:59.830695+09:30	2016-05-06	2016-05-07	2016-06-06	\N	0	2016-05-06 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
179	18	8	0	C	RollBed	rollbed fix every 3 days		2016-04-30 09:58:04.536759+09:30	2016-05-07	2016-05-08	2016-06-07	\N	29	2016-05-07 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
180	20	16	0	C	Printer	add paper to printer		2016-04-30 09:58:04.54672+09:30	2016-05-07	2016-05-08	2016-06-07	\N	0	2016-05-07 09:30:00+09:30	\N	f	\N	0.00	2.00	0.00	0.00	\N	\N
181	25	12	0	A	General Maintenance	do something every 4 days		2016-04-30 09:58:04.556646+09:30	2016-05-07	2016-05-08	2016-06-07	\N	30	2016-05-07 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
182	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n		2016-04-30 09:58:04.56617+09:30	2016-05-07	2016-05-08	2016-06-07	\N	0	2016-05-07 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
183	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-04-30 09:58:04.604326+09:30	2016-05-07	2016-05-08	2016-06-07	\N	0	2016-05-07 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
184	14	7	0	A	General Maintenance	first week maint		2016-05-01 00:58:08.700623+09:30	2016-05-02	2016-05-06	2016-06-02	\N	30	2016-05-02 00:00:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
185	21	12	0	A	General Maintenance	Tue maintenance		2016-05-01 00:58:08.710194+09:30	2016-05-02	2016-05-06	2016-06-02	\N	30	2016-05-02 00:00:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
186	19	6	0	C	Lube	relube the stud machine		2016-05-01 09:38:29.969401+09:30	2016-05-08	2016-05-09	2016-06-08	\N	30	2016-05-08 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
187	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n		2016-05-01 09:38:29.979623+09:30	2016-05-08	2016-05-09	2016-06-08	\N	0	2016-05-08 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
188	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-05-01 09:38:30.019623+09:30	2016-05-08	2016-05-09	2016-06-08	\N	0	2016-05-08 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
189	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n		2016-05-02 09:54:58.457284+09:30	2016-05-09	2016-05-10	2016-06-09	\N	0	2016-05-09 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
190	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-05-02 09:54:58.635747+09:30	2016-05-09	2016-05-10	2016-06-09	\N	0	2016-05-09 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
191	18	8	0	C	RollBed	rollbed fix every 3 days		2016-05-03 09:33:37.691813+09:30	2016-05-10	2016-05-11	2016-06-10	\N	29	2016-05-10 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
192	19	6	0	C	Lube	relube the stud machine		2016-05-03 09:33:37.745113+09:30	2016-05-10	2016-05-11	2016-06-10	\N	30	2016-05-10 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
193	20	16	0	C	Printer	add paper to printer		2016-05-03 09:33:37.754244+09:30	2016-05-10	2016-05-11	2016-06-10	\N	0	2016-05-10 09:30:00+09:30	\N	f	\N	0.00	2.00	0.00	0.00	\N	\N
194	28	29	0	A	General Maintenance	This one goes off every day with the following weapons\n		2016-05-03 09:33:37.767953+09:30	2016-05-10	2016-05-11	2016-06-10	\N	0	2016-05-10 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
195	29	28	0	A	General Maintenance	A new farm with these animals\n\n		2016-05-03 09:33:37.806548+09:30	2016-05-10	2016-05-11	2016-06-10	\N	0	2016-05-10 09:30:00+09:30	\N	f	\N	0.00	0.00	0.00	0.00	\N	\N
\.


--
-- Data for Name: task_check; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY task_check (task_id, seq, descr, done, done_date) FROM stdin;
143	1	the first check item	f	\N
143	2	second check item	f	\N
143	3	third item	f	\N
143	4	fourth item	f	\N
143	5	fifth item	f	\N
144	1	the first check item	f	\N
144	2	second check item	f	\N
144	3	third item	f	\N
144	4	fourth item	f	\N
144	5	fifth item	f	\N
145	1	Machine gun	f	\N
145	2	Rifle	f	\N
145	3	Chainsaw	f	\N
145	4	Atomic Bomb	f	\N
146	1	Machine gun	f	\N
146	2	Rifle	f	\N
146	3	Chainsaw	f	\N
146	4	Atomic Bomb	f	\N
147	1	Machine gun	f	\N
147	2	Rifle	f	\N
147	3	Chainsaw	f	\N
147	4	Atomic Bomb	f	\N
148	1	Machine gun	f	\N
148	2	Rifle	f	\N
148	3	Chainsaw	f	\N
148	4	Atomic Bomb	f	\N
149	1	Machine gun	f	\N
149	2	Rifle	f	\N
149	3	Chainsaw	f	\N
149	4	Atomic Bomb	f	\N
150	1	Machine gun	f	\N
150	2	Rifle	f	\N
150	3	Chainsaw	f	\N
150	4	Atomic Bomb	f	\N
152	1	horse	f	\N
152	2	cow	f	\N
152	3	duck	f	\N
152	4	turtle	f	\N
153	1	horse	f	\N
153	2	cow	f	\N
153	3	duck	f	\N
153	4	turtle	f	\N
154	1	horse	f	\N
154	2	cow	f	\N
154	3	duck	f	\N
154	4	turtle	f	\N
155	1	horse	f	\N
155	2	cow	f	\N
155	3	duck	f	\N
155	4	turtle	f	\N
156	1	horse	f	\N
156	2	cow	f	\N
156	3	duck	f	\N
156	4	turtle	f	\N
157	1	horse	f	\N
157	2	cow	f	\N
157	3	duck	f	\N
157	4	turtle	f	\N
159	1	horse	f	\N
159	2	cow	f	\N
159	3	duck	f	\N
159	4	turtle	f	\N
160	1	horse	f	\N
160	2	cow	f	\N
160	3	duck	f	\N
160	4	turtle	f	\N
161	1	horse	f	\N
161	2	cow	f	\N
161	3	duck	f	\N
161	4	turtle	f	\N
162	1	horse	f	\N
162	2	cow	f	\N
162	3	duck	f	\N
162	4	turtle	f	\N
163	1	horse	f	\N
163	2	cow	f	\N
163	3	duck	f	\N
163	4	turtle	f	\N
164	1	horse	f	\N
164	2	cow	f	\N
164	3	duck	f	\N
164	4	turtle	f	\N
165	2	cow	t	2016-04-25 23:28:21.565082+09:30
165	4	turtle	t	2016-04-25 23:30:12.018019+09:30
165	1	horse	t	2016-04-25 23:30:17.309883+09:30
165	3	duck	t	2016-04-25 23:30:58.195911+09:30
151	1	Machine gun	t	2016-04-26 00:11:07.347106+09:30
151	3	Chainsaw	t	2016-04-26 00:11:08.468363+09:30
151	2	Rifle	t	2016-04-26 00:11:10.902424+09:30
151	4	Atomic Bomb	t	2016-04-26 00:11:12.821116+09:30
158	2	cow	t	2016-04-26 00:15:01.601625+09:30
158	4	turtle	t	2016-04-26 00:15:02.549122+09:30
158	3	duck	t	2016-04-26 00:15:03.616334+09:30
158	1	horse	t	2016-04-26 00:15:04.554871+09:30
167	1	Machine gun	f	\N
167	2	Rifle	f	\N
167	3	Chainsaw	f	\N
167	4	Atomic Bomb	f	\N
168	1	horse	f	\N
168	2	cow	f	\N
168	3	duck	f	\N
168	4	turtle	f	\N
172	1	Machine gun	f	\N
172	2	Rifle	f	\N
172	3	Chainsaw	f	\N
172	4	Atomic Bomb	f	\N
173	1	horse	f	\N
173	2	cow	f	\N
173	3	duck	f	\N
173	4	turtle	f	\N
174	1	Machine gun	f	\N
174	2	Rifle	f	\N
174	3	Chainsaw	f	\N
174	4	Atomic Bomb	f	\N
175	1	horse	f	\N
175	2	cow	f	\N
175	3	duck	f	\N
175	4	turtle	f	\N
177	1	Machine gun	f	\N
177	2	Rifle	f	\N
177	3	Chainsaw	f	\N
177	4	Atomic Bomb	f	\N
178	1	horse	f	\N
178	2	cow	f	\N
178	3	duck	f	\N
178	4	turtle	f	\N
182	1	Machine gun	f	\N
182	2	Rifle	f	\N
182	3	Chainsaw	f	\N
182	4	Atomic Bomb	f	\N
183	1	horse	f	\N
183	2	cow	f	\N
183	3	duck	f	\N
183	4	turtle	f	\N
187	1	Machine gun	f	\N
187	2	Rifle	f	\N
187	3	Chainsaw	f	\N
187	4	Atomic Bomb	f	\N
188	1	horse	f	\N
188	2	cow	f	\N
188	3	duck	f	\N
188	4	turtle	f	\N
189	1	Machine gun	f	\N
189	2	Rifle	f	\N
189	3	Chainsaw	f	\N
189	4	Atomic Bomb	f	\N
190	1	horse	f	\N
190	2	cow	f	\N
190	3	duck	f	\N
190	4	turtle	f	\N
194	1	Machine gun	f	\N
194	2	Rifle	f	\N
194	3	Chainsaw	f	\N
194	4	Atomic Bomb	f	\N
195	1	horse	f	\N
195	2	cow	f	\N
195	3	duck	f	\N
195	4	turtle	f	\N
\.


--
-- Name: task_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('task_id_seq', 195, true);


--
-- Data for Name: task_part; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY task_part (task_id, part_id, qty, notes) FROM stdin;
143	72	2.00	2 blade items
143	58	3.00	3 block itiems
143	62	3.00	some die bushes as well
144	72	2.00	2 blade items
144	58	3.00	3 block itiems
144	62	3.00	some die bushes as well
145	68	22.00	calibre
145	64	66.00	die steel
146	68	22.00	calibre
146	64	66.00	die steel
147	68	22.00	calibre
147	64	66.00	die steel
148	68	22.00	calibre
148	64	66.00	die steel
149	68	22.00	calibre
149	64	66.00	die steel
150	68	22.00	calibre
150	64	66.00	die steel
151	68	22.00	calibre
151	64	66.00	die steel
167	68	22.00	calibre
167	64	66.00	die steel
172	68	22.00	calibre
172	64	66.00	die steel
174	68	22.00	calibre
174	64	66.00	die steel
177	68	22.00	calibre
177	64	66.00	die steel
182	68	22.00	calibre
182	64	66.00	die steel
187	68	22.00	calibre
187	64	66.00	die steel
189	68	22.00	calibre
189	64	66.00	die steel
194	68	22.00	calibre
194	64	66.00	die steel
\.


--
-- Data for Name: user_log; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY user_log (id, logdate, duration, ms, func, input, output) FROM stdin;
\.


--
-- Name: user_log_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('user_log_id_seq', 19046, true);


--
-- Data for Name: user_role; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY user_role (user_id, site_id, worker, sitemgr, contractor) FROM stdin;
\.


--
-- Data for Name: user_site; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY user_site (user_id, site_id, role) FROM stdin;
20	2	Admin
8	2	Site Manager
21	2	Floor
22	2	Admin
23	2	Admin
24	2	Floor
25	1	Site Manager
26	1	Worker
27	1	Worker
28	1	Worker
29	1	Worker
30	2	Site Manager
31	2	Worker
32	5	Site Manager
33	2	Admin
35	6	Site Manager
36	6	Worker
37	7	Site Manager
38	7	Worker
39	8	Site Manager
40	9	Worker
8	4	Site Manager
8	1	Site Manager
48	7	Worker
49	4	Worker
44	2	Site Manager
44	4	Site Manager
46	2	Site Manager
47	7	Worker
45	2	Worker
45	4	Worker
50	6	Worker
1	7	Admin
1	13	Admin
1	1	Admin
43	12	Service Contractor
41	1	Service Contractor
41	3	Service Contractor
34	8	Worker
34	3	Worker
34	4	Worker
34	6	Worker
34	7	Worker
25	2	
35	2	Site Manager
33	7	Admin
33	8	Admin
30	3	Site Manager
29	2	Worker
41	2	Service Contractor
28	2	Worker
43	2	Service Contractor
27	2	Worker
48	2	Worker
49	2	Worker
54	2	Public
50	2	Worker
1	3	Admin
1	4	Admin
1	2	Admin
\.


--
-- Data for Name: user_skill; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY user_skill (user_id, skill_id) FROM stdin;
20	1
8	1
48	5
22	4
22	5
23	1
25	2
26	3
27	3
28	3
29	4
30	5
31	2
32	5
34	2
35	5
36	3
37	5
38	3
39	5
40	3
49	5
50	5
44	4
1	1
1	4
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY users (id, username, passwd, name, address, email, sms, site_id, role, notes) FROM stdin;
20	pzitis	passwd	Peter Zitis		pzitis@email.com	TBA ......	2	Admin	
21	Edinburgh-Main	passwd	Edinburgh Main Factory		Edinburgh-Main@email.com		2	Floor	
22	webdev	passwd	Web Development Team Upstairs		webdev@email.com		2	Admin	
24	floor	passwd	Floor Control demo		floor@email.com		2	Floor	
25	wayne.cridland	passwd	Wayne Cridland		wayne.cridland@email.com		1	Site Manager	
30	stephen.green	passwd	Stephen Green		stephen.green@email.com		2	Site Manager	
32	roger.giles	passwd	Roger Giles		roger.giles@email.com		5	Site Manager	
35	dan.bright	passwd	Dan Bright		dan.bright@email.com		6	Site Manager	
37	richie.maz	passwd	Richie Mazurkiewicz		richie.maz@email.com		7	Site Manager	
44	sm1	passwd	Sm1		sm1@email.com		2	Site Manager	
39	john.klomp	passwd	John Klomp		john.klomp@email.com		8	Site Manager	
42	matt.keane	passwd	Matt Keane		matt.keane@email.com		0	Service Contractor	
43	ryan.merret	passwd	Ryan Merret		ryan.merret@email.com		0	Service Contractor	
8	shane.voigt	passwd	Shane Voigt		shane.voigt@email.com	Unknown ..	2	Site Manager	
46	testsm1	passwd	Testsm1		testsm1@email.com		2	Site Manager	
1	admin	passwd	Admin Bootstrap User		admin@email.com		1	Admin	<h2>Notes</h2><br>Add some notes to the admin user
54	thy thise	passwd	new name		thy thise@email.com	234234234234	0	Public	
41	kieth.morton	passwd	Kieth Morton		kieth.morton@email.com	34234	0	Service Contractor	
28	ralph.iengo	passwd	Ralph Iengo		ralph.iengo@email.com		1	Technician	
29	brett.robson	passwd	Brett Robson		brett.robson@email.com		1	Technician	
31	matt.howell	passwd	Matt Howell		matt.howell@email.com		2	Technician	
34	jim.bratis	passwd	Jim Bratis		jim.bratis@email.com		5	Technician	
36	lee.reed	passwd	Lee Reed		lee.reed@email.com		6	Technician	
23	guest aa	passwd	Guest admin user		guest aa@email.com		2	Admin	
38	sachin.sharma	passwd	Sachin Sharma		sachin.sharma@email.com		7	Technician	
40	scott.arris	passwd	Scott Arris		scott.arris@email.com		9	Technician	
48	testwminto	passwd	Testwminto		testwminto@email.com		7	Technician	
49	testsmic	passwd	Testsmic		testsmic@email.com		4	Technician	
47	mintow1	passwd	Mintow1		mintow1@email.com		7	Technician	
50	tomagotest	passwd	Tomagotest		tomagotest@email.com		6	Technician	
27	sam.combes	passwd	Sam Combes		sam.combes@email.com		1	Technician	
26	ashley.stroeger	passwd	Ashley Stroeger		ashley.stroeger@email.com		1	Technician	
45	testw1	passwd	Test Worker 1		testw1@email.com	0415083977	2	Technician	
33	postgres	passwd	postgres Oconnor		postgres@email.com		2	Admin	
\.


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('users_id_seq', 54, true);


--
-- Data for Name: vendor; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY vendor (id, name, descr, address, phone, fax, contact_name, contact_email, orders_email, rating, notes) FROM stdin;
2	IAS	Kieth x	40 Barfield cres north	0444 Unknown		Keith Morton	kmorton@ias.com.au		Great	
\.


--
-- Name: vendor_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('vendor_id_seq', 3, true);


--
-- Data for Name: vendor_price; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY vendor_price (part_id, vendor_id, datefrom, price, min_qty, notes) FROM stdin;
90	2	2015-12-09 13:43:23.250349+10:30	889.00	1.00	
86	2	2015-12-09 13:44:58.79913+10:30	499.00	1.00	
113	2	2015-12-09 13:47:15.634444+10:30	480.00	1.00	
85	2	2015-12-09 13:47:15.637223+10:30	397.00	1.00	
87	2	2015-12-09 13:47:15.639259+10:30	328.00	12.00	
\.


--
-- Data for Name: wo_assignee; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY wo_assignee (id, user_id) FROM stdin;
9	1
9	26
9	29
10	1
10	26
10	29
10	35
10	21
11	35
11	37
12	24
12	23
12	34
12	29
12	1
12	26
12	20
13	30
13	41
13	39
13	42
13	33
14	30
14	41
14	39
14	42
14	33
15	1
16	33
17	33
18	33
19	33
20	33
21	33
22	1
23	33
23	1
23	26
24	1
25	1
26	1
27	1
28	1
29	1
30	1
31	1
32	1
33	1
34	1
35	1
36	1
37	1
38	29
39	35
41	1
42	1
43	35
46	35
\.


--
-- Data for Name: wo_docs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY wo_docs (id, doc_id) FROM stdin;
16	31
16	22
16	10
16	28
16	29
16	30
17	31
17	22
18	35
18	37
18	36
19	31
19	29
20	31
20	22
20	10
20	28
21	31
21	22
21	28
22	34
22	32
22	54
22	55
22	56
22	57
22	58
23	34
23	32
23	59
23	60
23	54
24	32
24	33
24	61
25	34
25	61
25	62
26	30
26	22
26	10
26	46
26	47
26	48
26	49
26	50
26	64
27	27
27	66
32	34
32	65
33	34
33	65
34	27
34	66
35	34
35	32
35	65
35	72
36	32
36	34
36	65
36	72
37	34
37	53
38	10
38	46
38	47
38	48
38	49
38	50
38	64
39	33
39	61
39	62
42	27
42	66
\.


--
-- Data for Name: wo_skills; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY wo_skills (id, skill_id) FROM stdin;
9	4
9	6
9	1
10	4
10	6
10	1
10	5
10	2
10	3
11	6
11	4
11	3
12	6
12	3
13	5
13	1
13	4
13	2
13	3
13	6
14	5
14	1
14	4
14	2
14	3
14	6
15	4
16	4
17	4
18	6
19	6
19	5
20	6
20	1
20	3
21	6
21	1
21	3
23	6
23	1
24	3
25	4
26	3
27	3
28	4
29	6
30	4
31	5
32	3
33	3
34	4
35	4
36	6
37	4
38	6
39	5
41	6
42	2
43	2
46	6
\.


--
-- Data for Name: workorder; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY workorder (id, event_id, startdate, est_duration, actual_duration, descr, status, notes) FROM stdin;
2	210	2016-01-12 12:00:00+10:30	60	0	create workorder test 1	Assigned	<h2>Site: Edinburgh - Factory</h2>\n40 Barfield Crescent, Edinburgh North, SA 5113<p>\n<h2>Machine: Bracket B4/5</h2>\nSome notes here for the Bracket Machine<br><br><ol><li>Do this</li><li>Do that</li><li>Do the other</li><li>Do some more</li></ol><p>\n<h2>Tool: Guillo</h2>\nMaintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol><p>\neNotes.SiteNotes<p>\n
3	210	2016-01-12 12:00:00+10:30	90	0	workorder test 2	Assigned	<h2>Site: Edinburgh - Factory</h2>\n40 Barfield Crescent, Edinburgh North, SA 5113<p>\n<h2>Machine: Bracket B4/5</h2>\nSome notes here for the Bracket Machine<br><br><ol><li>Do this</li><li>Do that</li><li>Do the other</li><li>Do some more</li></ol><p>\n<h2>Tool: Guillo</h2>\nMaintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol><p>\neNotes.SiteNotes<p>\n
4	210	2016-01-12 12:00:00+10:30	90	0	workorder test 2	Assigned	<h2>Site: Edinburgh - Factory</h2>\n40 Barfield Crescent, Edinburgh North, SA 5113<p>\n<h2>Machine: Bracket B4/5</h2>\nSome notes here for the Bracket Machine<br><br><ol><li>Do this</li><li>Do that</li><li>Do the other</li><li>Do some more</li></ol><p>\n<h2>Tool: Guillo</h2>\nMaintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol><p>\neNotes.SiteNotes<p>\n
5	210	2016-01-12 12:00:00+10:30	90	0	workorder test 3	Assigned	<h2>Site: Edinburgh - Factory</h2>\n40 Barfield Crescent, Edinburgh North, SA 5113<p>\neNotes.SiteNotes<p>\n<h2>Machine: Bracket B4/5</h2>\nSome notes here for the Bracket Machine<br><br><ol><li>Do this</li><li>Do that</li><li>Do the other</li><li>Do some more</li></ol><p>\n<h2>Tool: Guillo</h2>\nMaintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol><p>\n
6	210	2016-01-12 12:00:00+10:30	90	0	workorder test 4	Assigned	<h2>Site: Edinburgh - Factory</h2>\n40 Barfield Crescent, Edinburgh North, SA 5113<p>\neNotes.SiteNotes<p>\n<h2>Machine: Bracket B4/5</h2>\nSome notes here for the Bracket Machine<br><br><ol><li>Do this</li><li>Do that</li><li>Do the other</li><li>Do some more</li></ol><p>\n<h2>Tool: Guillo</h2>\nMaintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol><p>\n
7	210	2016-01-12 12:00:00+10:30	90	0	workorder test 5	Assigned	<h2>Site: Edinburgh - Factory</h2>\n40 Barfield Crescent, Edinburgh North, SA 5113<p>\n<h1>Main Edinburgh Site</h1><div><br></div><div>This includes all machines at the main factory, and any machines at sub-factories at the main Edinburgh site</div><p>\n<h2>Machine: Bracket B4/5</h2>\nSome notes here for the Bracket Machine<br><br><ol><li>Do this</li><li>Do that</li><li>Do the other</li><li>Do some more</li></ol><p>\n<h2>Tool: Guillo</h2>\nMaintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol><p>\n
8	210	2016-01-12 12:00:00+10:30	90	0	workorder test 6	Assigned	<h2>Site: Edinburgh - Factory</h2>\n40 Barfield Crescent, Edinburgh North, SA 5113<p>\n<h1>Main Edinburgh Site</h1><div><br></div><div>This includes all machines at the main factory, and any machines at sub-factories at the main Edinburgh site</div><p>\n<h2>Machine: Bracket B4/5</h2>\nSome notes here for the Bracket Machine<br><br><ol><li>Do this</li><li>Do that</li><li>Do the other</li><li>Do some more</li></ol><p>\n<h2>Tool: Guillo</h2>\nMaintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol><p>\n
9	210	2016-01-12 12:00:00+10:30	90	0	workorder test 7	Assigned	<h2>Site: Edinburgh - Factory</h2>\n40 Barfield Crescent, Edinburgh North, SA 5113<p>\n<h1>Main Edinburgh Site</h1><div><br></div><div>This includes all machines at the main factory, and any machines at sub-factories at the main Edinburgh site</div><p>\n<h2>Machine: Bracket B4/5</h2>\nSome notes here for the Bracket Machine<br><br><ol><li>Do this</li><li>Do that</li><li>Do the other</li><li>Do some more</li></ol><p>\n<h2>Tool: Guillo</h2>\nMaintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol><p>\n
10	210	2016-01-23 15:00:00+10:30	150	0	scheduled fix for this machine - checx all the notes and things	Assigned	<h2>Site: Edinburgh - Factory</h2>\n40 Barfield Crescent, Edinburgh North, SA 5113<p>\n<h1>Main Edinburgh Site</h1><div><br></div><div>This includes all machines at the main factory, and any machines at sub-factories at the main Edinburgh site</div><p>\n<h2>Machine: Bracket B4/5</h2>\nSome notes here for the Bracket Machine<br><br><ol><li>Do this</li><li>Do that</li><li>Do the other</li><li>Do some more</li></ol><p>\n<h2>Tool: Guillo</h2>\nMaintenance instructions for the Guillo<br><br><ol><li>do the initial safety steps</li><li>do the other steps</li><li>do the final steps</li></ol><p>\n
11	209	2016-01-25 16:00:00+10:30	480	0	sched maintenance - strip down and rebuild the machine	Assigned	
12	208	2016-02-01 16:00:00+10:30	240	0	sched maintenace - chain relube	Assigned	
13	210	2016-01-23 16:00:00+10:30	600	0	Regular Preventative Maint and inspection	Assigned	
14	210	2016-01-23 16:00:00+10:30	600	0	Regular Preventative Maint and inspection	Assigned	
15	210	2016-01-14 10:00:00+10:30	30	0	test case - doc support	Assigned	
16	210	2016-01-14 10:00:00+10:30	60	0	test case - should have a number of docs attached	Assigned	
17	210	2016-01-14 10:00:00+10:30	60	0	test case - should have a number of docs attached	Assigned	
18	208	2016-01-14 11:00:00+10:30	30	0	test case, simple docs attached	Assigned	
19	210	2016-01-14 14:00:00+10:30	30	0	test case - tx through sbs mail server	Assigned	
20	210	2016-01-15 09:00:00+10:30	90	0	test of large attachments through xchange server	Assigned	
21	210	2016-01-15 09:00:00+10:30	90	0	test of large attachments through xchange server	Assigned	
22	233	2016-01-31 11:32:39.677+10:30	90	0	oe	Assigned	
23	233	2016-02-07 14:33:50.175+10:30	90	0	2pm on feb 7	Assigned	
24	234	2016-02-01 16:09:59.009+10:30	120	0	please fix the borked machine	Assigned	
25	234	2016-01-31 10:12:50.315+10:30	90	0	fix this item	Assigned	
26	229	2016-01-30 11:52:35.816+10:30	30	0	todo	Assigned	
27	236	2016-01-31 12:55:29.3+10:30	0	0	can you sharpen the blade please	Assigned	
28	235	2016-02-05 08:51:07.262+10:30	90	0	do the work	Assigned	
29	235	2016-02-06 09:15:40.922+10:30	60	0	try again - this time with notes attached	Assigned	
30	235	2016-02-05 09:18:28.382+10:30	120	0	where are the notes for this workorder ?	Assigned	
31	235	2016-02-27 09:22:39.188+10:30	90	0	storing notes now on backend	Assigned	these notes should appear in the email
32	235	2016-02-12 09:26:03.602+10:30	90	0	needs fixing	Assigned	these notes to appear just under the estimated duration
33	235	2016-02-12 09:26:03.602+10:30	90	0	needs replacing	Assigned	these notes to appear just under the estimated duration<br><br>consider these points :<br><ul><li>this is a list that appears inside the other list</li><li>and so is this item</li></ul><div><br></div><div><ol><li>and a numbered list as well</li><li>with a few lines on there</li></ol></div>
36	235	2016-02-04 16:00:00+10:30	120	0	4th at 4pm	Assigned	some notes here
34	236	2016-02-06 10:42:17.437+10:30	150	0	testing sending time data	Assigned	time was set to 12, but should be 10:44 in the startdate field
35	235	2016-02-12 11:00:00+10:30	90	0	7pm	Assigned	fix all issues
37	232	2016-02-13 13:00:00+10:30	300	0	The 13th at 13 oclock	Assigned	should be lucky this time<br><br>
38	229	2016-02-13 13:00:00+10:30	90	0	another 13th at 13 oclock	Assigned	another lucky one for sure
39	234	2016-02-02 14:00:00+10:30	150	0	work to be done	Assigned	date should read 2pm today
40	237	2016-02-12 14:00:00+10:30	0	0		Assigned	
41	237	2016-02-26 15:00:00+10:30	90	0	a test with 2 email outputs	Assigned	some notes
42	236	2016-02-27 11:00:00+10:30	90	0	one more	Assigned	<ol><li>a list</li><li>another line in the same list</li></ol>
43	237	2016-02-14 11:00:00+10:30	150	0	test of updating the notes field	Assigned	initial notes
44	234	2016-03-01 10:00:00+10:30	90	0	test of adding workorder, needs to update workorder list	Assigned	
45	234	2016-02-08 10:00:00+10:30	30	0	should add a sixth line to the workorder list	Assigned	
46	265	2016-03-04 16:00:00+10:30	90	0	fix this one	Assigned	
\.


--
-- Name: workorder_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('workorder_id_seq', 46, true);


--
-- Name: doc_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY doc
    ADD CONSTRAINT doc_pkey PRIMARY KEY (id);


--
-- Name: doc_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY doc_type
    ADD CONSTRAINT doc_type_pkey PRIMARY KEY (id);


--
-- Name: event_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY event
    ADD CONSTRAINT event_pkey PRIMARY KEY (id);


--
-- Name: event_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY event_type
    ADD CONSTRAINT event_type_pkey PRIMARY KEY (id);


--
-- Name: hashtag_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY hashtag
    ADD CONSTRAINT hashtag_pkey PRIMARY KEY (id);


--
-- Name: machine_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY machine
    ADD CONSTRAINT machine_pkey PRIMARY KEY (id);


--
-- Name: part_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY part
    ADD CONSTRAINT part_pkey PRIMARY KEY (id);


--
-- Name: sched_control_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY sched_control
    ADD CONSTRAINT sched_control_pkey PRIMARY KEY (id);


--
-- Name: sched_task_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY sched_task
    ADD CONSTRAINT sched_task_pkey PRIMARY KEY (id);


--
-- Name: site_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY site
    ADD CONSTRAINT site_pkey PRIMARY KEY (id);


--
-- Name: skill_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY skill
    ADD CONSTRAINT skill_pkey PRIMARY KEY (id);


--
-- Name: sm_task_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY sm_task
    ADD CONSTRAINT sm_task_pkey PRIMARY KEY (id);


--
-- Name: stock_level_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY stock_level
    ADD CONSTRAINT stock_level_pkey PRIMARY KEY (part_id);


--
-- Name: sys_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY sys_log
    ADD CONSTRAINT sys_log_pkey PRIMARY KEY (id);


--
-- Name: task_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY task
    ADD CONSTRAINT task_pkey PRIMARY KEY (id);


--
-- Name: user_log_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY user_log
    ADD CONSTRAINT user_log_pkey PRIMARY KEY (id);


--
-- Name: users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: vendor_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY vendor
    ADD CONSTRAINT vendor_pkey PRIMARY KEY (id);


--
-- Name: workorder_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace: 
--

ALTER TABLE ONLY workorder
    ADD CONSTRAINT workorder_pkey PRIMARY KEY (id);


--
-- Name: component_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX component_idx ON component USING btree (machine_id, id);


--
-- Name: component_part_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX component_part_idx ON component_part USING btree (component_id, part_id);


--
-- Name: component_position_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX component_position_idx ON component USING btree (machine_id, "position");


--
-- Name: doc_path_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX doc_path_idx ON doc USING btree (path);


--
-- Name: doc_rev_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX doc_rev_idx ON doc_rev USING btree (doc_id, id);


--
-- Name: event_allocation_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX event_allocation_idx ON event USING btree (allocated_to, id);


--
-- Name: event_doc_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX event_doc_idx ON event_doc USING btree (event_id, doc_id);


--
-- Name: event_site_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX event_site_idx ON event USING btree (site_id, startdate);


--
-- Name: part_price_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX part_price_idx ON part_price USING btree (part_id, datefrom);


--
-- Name: part_stock_code_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX part_stock_code_idx ON part USING btree (stock_code);


--
-- Name: part_stock_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX part_stock_idx ON part_stock USING btree (part_id, datefrom);


--
-- Name: part_vendor_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX part_vendor_idx ON part_vendor USING btree (part_id, vendor_id);


--
-- Name: sched_task_part_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX sched_task_part_idx ON sched_task_part USING btree (task_id, part_id);


--
-- Name: site_layout_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX site_layout_idx ON site_layout USING btree (site_id, seq);


--
-- Name: sm_component_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX sm_component_idx ON sm_tool USING btree (task_id, machine_id, tool_id);


--
-- Name: sm_component_item_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX sm_component_item_idx ON sm_component_item USING btree (task_id, component, seq);


--
-- Name: sm_machine_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX sm_machine_idx ON sm_machine USING btree (task_id, machine_id);


--
-- Name: sm_machine_item_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX sm_machine_item_idx ON sm_machine_item USING btree (task_id, machine_id, seq);


--
-- Name: sm_part_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX sm_part_idx ON sm_parts USING btree (task_id, part_id, date);


--
-- Name: sm_task_user_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX sm_task_user_idx ON sm_task USING btree (user_id, date);


--
-- Name: sm_tool_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX sm_tool_idx ON sm_tool USING btree (task_id, machine_id, tool_id);


--
-- Name: sm_tool_item_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX sm_tool_item_idx ON sm_tool_item USING btree (task_id, tool_id, seq);


--
-- Name: stock_level_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE INDEX stock_level_idx ON stock_level USING btree (part_id, site_id);


--
-- Name: sys_log_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX sys_log_idx ON sys_log USING btree (logdate, id);


--
-- Name: task_check_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX task_check_idx ON task_check USING btree (task_id, seq);


--
-- Name: task_part_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX task_part_idx ON task_part USING btree (task_id, part_id);


--
-- Name: user_role_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX user_role_idx ON user_role USING btree (user_id, site_id);


--
-- Name: user_site_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX user_site_idx ON user_site USING btree (user_id, site_id);


--
-- Name: user_skill_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX user_skill_idx ON user_skill USING btree (user_id, skill_id);


--
-- Name: vendor_price_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX vendor_price_idx ON vendor_price USING btree (part_id, vendor_id, datefrom);


--
-- Name: wo_assignee_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX wo_assignee_idx ON wo_assignee USING btree (id, user_id);


--
-- Name: wo_docs_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX wo_docs_idx ON wo_docs USING btree (id, doc_id);


--
-- Name: wo_skills_idx; Type: INDEX; Schema: public; Owner: postgres; Tablespace: 
--

CREATE UNIQUE INDEX wo_skills_idx ON wo_skills USING btree (id, skill_id);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

