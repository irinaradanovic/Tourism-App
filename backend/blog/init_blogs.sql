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

INSERT INTO blogs (id, author_id, title, description, images, created_at) VALUES
('blog_1', '2', 'Letovanje u Grčkoj 2026', 'Sjajno iskustvo na Sitoniji, prelepe plaže i čisto more.', ARRAY['uploads/blogs/image1.jpg', 'uploads/blogs/image2.jpg'], NOW() - INTERVAL '5 days'),
('blog_2', '3', 'Planinarenje na Tari', 'Vikend avantura i pogled sa Banjske stene koji ostavlja bez daha.', ARRAY['uploads/blogs/image3.jpg'], NOW() - INTERVAL '4 days'),
('blog_3', '4', 'Vikend u Budimpešti', 'Vodič kroz najbolja mesta za hranu i obilazak parlamenta.', ARRAY['uploads/blogs/image4.jpg', 'uploads/blogs/image5.jpg', 'uploads/blogs/image6.jpeg'], NOW() - INTERVAL '3 days'),
('blog_4', '5', 'Skijanje na Kopaoniku', 'Staze su odlično pripremljene ove sezone, gužve su umerene.', ARRAY['uploads/blogs/image2.jpg', 'uploads/blogs/image4.jpg'], NOW() - INTERVAL '2 days'),
('blog_5', '6', 'Istraživanje Rima', 'Koloseum, Vatikan i nezaobilazni italijanski sladoled.', ARRAY['uploads/blogs/image1.jpg', 'uploads/blogs/image3.jpg', 'uploads/blogs/image5.jpg', 'uploads/blogs/image6.jpeg'], NOW() - INTERVAL '1 day'),
('blog_6', '7', 'Zlatibor sa decom', 'Najbolje rute za šetnju i aktivnosti za najmlađe.', ARRAY['uploads/blogs/image2.jpg'], NOW()),
('blog_7', '8', 'Jednodnevni izlet u Novi Sad', 'Šetnja Dunavskom ulicom i zalazak sunca na Petrovaradinskoj tvrđavi.', ARRAY['uploads/blogs/image3.jpg', 'uploads/blogs/image1.jpg'], NOW()),
('blog_8', '12', 'Lepote Tare u proleće', 'Smaragdno zelena Drina i mirisi šume.', ARRAY['uploads/blogs/image5.jpg'], NOW()),
('blog_9', '15', 'Kulinarski vodič kroz Sarajevo', 'Gde pojesti najbolje ćevape i probati autentičnu bosansku kafu.', ARRAY['uploads/blogs/image6.jpeg', 'uploads/blogs/image4.jpg', 'uploads/blogs/image2.jpg'], NOW());

INSERT INTO likes (user_id, blog_id) VALUES
('2', 'blog_2'),
('2', 'blog_3'),
('3', 'blog_1'),
('3', 'blog_4'),
('4', 'blog_1'),
('4', 'blog_5'),
('5', 'blog_2'),
('5', 'blog_3'),
('6', 'blog_1'),
('7', 'blog_5'),
('8', 'blog_6'),
('9', 'blog_4'),
('10', 'blog_1'),
('11', 'blog_2'),
('12', 'blog_7'),
('13', 'blog_3'),
('14', 'blog_9'),
('15', 'blog_8'),
('16', 'blog_1'),
('16', 'blog_9');