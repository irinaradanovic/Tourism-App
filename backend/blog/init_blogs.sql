CREATE TABLE IF NOT EXISTS blogs (
    id VARCHAR(255) PRIMARY KEY,
    author_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    images TEXT[], 
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS likes (
    id SERIAL PRIMARY KEY, 
    user_id VARCHAR(255) NOT NULL,
    blog_id VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_blog ON likes (user_id, blog_id);

-- Obriši stare podatke ako ponovo pokrećeš skriptu
TRUNCATE TABLE likes RESTART IDENTITY CASCADE;
TRUNCATE TABLE blogs RESTART IDENTITY CASCADE;

INSERT INTO blogs (id, author_id, title, description, images, created_at) VALUES
('blog_1', '2', 'Letovanje u Grčkoj 2026', 'Sjajno iskustvo na Sitoniji, prelepe plaže i čisto more.', ARRAY['uploads/blogs/image1.jpg', 'uploads/blogs/image2.jpg'], NOW() - INTERVAL '15 days'),
('blog_2', '3', 'Planinarenje na Tari', 'Vikend avantura i pogled sa Banjske stene koji ostavlja bez daha.', ARRAY['uploads/blogs/image3.jpg'], NOW() - INTERVAL '14 days'),
('blog_3', '4', 'Vikend u Budimpešti', 'Vodič kroz najbolja mesta za hranu i obilazak parlamenta.', ARRAY['uploads/blogs/image4.jpg', 'uploads/blogs/image5.jpg'], NOW() - INTERVAL '13 days'),
('blog_4', '5', 'Skijanje na Kopaoniku', 'Staze su odlično pripremljene ove sezone, gužve su umerene.', ARRAY['uploads/blogs/image2.jpg', 'uploads/blogs/image4.jpg'], NOW() - INTERVAL '12 days'),
('blog_5', '6', 'Istraživanje Rima', 'Koloseum, Vatikan i nezaobilazni italijanski sladoled.', ARRAY['uploads/blogs/image1.jpg', 'uploads/blogs/image3.jpg'], NOW() - INTERVAL '11 days'),
('blog_6', '7', 'Zlatibor sa decom', 'Najbolje rute za šetnju i aktivnosti za najmlađe.', ARRAY['uploads/blogs/image2.jpg'], NOW() - INTERVAL '10 days'),
('blog_7', '8', 'Jednodnevni izlet u Novi Sad', 'Šetnja Dunavskom ulicom i zalazak sunca na Petrovaradinskoj tvrđavi.', ARRAY['uploads/blogs/image3.jpg', 'uploads/blogs/image1.jpg'], NOW() - INTERVAL '9 days'),
('blog_8', '12', 'Lepote Tare u proleće', 'Smaragdno zelena Drina i mirisi šume.', ARRAY['uploads/blogs/image5.jpg'], NOW() - INTERVAL '8 days'),
('blog_9', '15', 'Kulinarski vodič kroz Sarajevo', 'Gde pojesti najbolje ćevape i probati autentičnu bosansku kafu.', ARRAY['uploads/blogs/image6.jpeg', 'uploads/blogs/image4.jpg'], NOW() - INTERVAL '7 days'),
('blog_10', '17', 'Čari Praga u proleće', 'Šetnja Karlovim mostom i degustacija čuvenog trdelnika.', ARRAY['uploads/blogs/image1.jpg'], NOW() - INTERVAL '6 days'),
('blog_11', '19', 'Zeleni pejzaži Slovenije', 'Bledsko jezero i mistična Postojnska jama su nas oduševili.', ARRAY['uploads/blogs/image3.jpg', 'uploads/blogs/image4.jpg'], NOW() - INTERVAL '5 days'),
('blog_12', '21', 'Adrenalinski vikend u kanjonu Tare', 'Rafting divljom rekom i kampovanje pod zvezdama.', ARRAY['uploads/blogs/image2.jpg'], NOW() - INTERVAL '4 days'),
('blog_13', '23', 'Barselona i arhitektura Gaudija', 'Fascinantna Sagrada Familia i opuštanje na plaži Barseloneta.', ARRAY['uploads/blogs/image5.jpg', 'uploads/blogs/image6.jpeg'], NOW() - INTERVAL '3 days'),
('blog_14', '25', 'Skriveni dragulji Ohridskog jezera', 'Obilazak crkve Sveti Jovan Kaneo i vožnja čamcem kroz izvore.', ARRAY['uploads/blogs/image1.jpg'], NOW() - INTERVAL '2 days'),
('blog_15', '26', 'Vikend u Beču', 'Dvorac Šenbrun, zimska palata i bečka šnicla u srcu grada.', ARRAY['uploads/blogs/image4.jpg'], NOW() - INTERVAL '1 day'),
('blog_16', '18', 'Krit - rajske plaže Grčke', 'Laguna Balos i plaža Elafonisi sa ružičastim peskom.', ARRAY['uploads/blogs/image2.jpg', 'uploads/blogs/image3.jpg'], NOW()),
('blog_17', '20', 'Istanbul na dva kontinenta', 'Krstarenje Bosforom, poseta Aja Sofiji i cenkanje na Kapali čaršiji.', ARRAY['uploads/blogs/image5.jpg'], NOW()),
('blog_18', '24', 'Etno sela Srbije', 'Mir, tišina i vrhunska domaća hrana u okolini Gornjeg Milanovca.', ARRAY['uploads/blogs/image1.jpg', 'uploads/blogs/image6.jpeg'], NOW());

INSERT INTO likes (user_id, blog_id) VALUES
('2', 'blog_2'), ('2', 'blog_3'), ('2', 'blog_10'),
('3', 'blog_1'), ('3', 'blog_4'), ('3', 'blog_11'),
('4', 'blog_1'), ('4', 'blog_5'), ('4', 'blog_12'),
('5', 'blog_2'), ('5', 'blog_3'), ('5', 'blog_13'),
('6', 'blog_1'), ('6', 'blog_14'),
('7', 'blog_5'), ('7', 'blog_15'),
('8', 'blog_6'), ('8', 'blog_16'),
('9', 'blog_4'), ('9', 'blog_17'),
('10', 'blog_1'), ('10', 'blog_18'),
('11', 'blog_2'), ('11', 'blog_10'),
('12', 'blog_7'), ('12', 'blog_11'),
('13', 'blog_3'), ('13', 'blog_12'),
('14', 'blog_9'), ('14', 'blog_13'),
('15', 'blog_8'), ('15', 'blog_14'),
('16', 'blog_1'), ('16', 'blog_9'), ('16', 'blog_15'),
('17', 'blog_11'), ('18', 'blog_12'), ('19', 'blog_13'),
('20', 'blog_14'), ('21', 'blog_15'), ('22', 'blog_16'),
('23', 'blog_17'), ('24', 'blog_18'), ('25', 'blog_1'),
('26', 'blog_2'), ('26', 'blog_9'), ('22', 'blog_1');