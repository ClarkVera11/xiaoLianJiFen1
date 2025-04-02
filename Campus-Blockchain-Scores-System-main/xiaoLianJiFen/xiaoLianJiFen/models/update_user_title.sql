DELIMITER //

CREATE TRIGGER update_user_title
BEFORE UPDATE ON users
FOR EACH ROW
BEGIN
    IF NEW.points != OLD.points THEN
        CASE
            WHEN NEW.points >= 400 THEN
                SET NEW.title = '荣耀王者';
            WHEN NEW.points >= 350 THEN
                SET NEW.title = '最强王者';
            WHEN NEW.points >= 300 THEN
                SET NEW.title = '至尊星耀';
            WHEN NEW.points >= 250 THEN
                SET NEW.title = '永恒钻石';
            WHEN NEW.points >= 200 THEN
                SET NEW.title = '尊贵铂金';
            WHEN NEW.points >= 150 THEN
                SET NEW.title = '荣耀黄金';
            WHEN NEW.points >= 100 THEN
                SET NEW.title = '秩序白银';
            ELSE
                SET NEW.title = '倔强青铜';
        END CASE;
    END IF;
END //

DELIMITER ;