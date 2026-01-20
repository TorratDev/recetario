-- Insert sample ingredients
INSERT INTO ingredients (name, category) VALUES
('flour', 'pantry'),
('sugar', 'pantry'),
('eggs', 'dairy'),
('milk', 'dairy'),
('butter', 'dairy'),
('olive oil', 'oils'),
('salt', 'seasonings'),
('pepper', 'seasonings'),
('garlic', 'vegetables'),
('onions', 'vegetables'),
('tomatoes', 'vegetables'),
('basil', 'herbs'),
('oregano', 'herbs'),
('chicken breast', 'meat'),
('ground beef', 'meat'),
('spaghetti', 'pasta'),
('rice', 'grains'),
('parmesan cheese', 'dairy'),
('mozzarella', 'dairy'),
('bell peppers', 'vegetables');

-- Insert sample tags
INSERT INTO tags (name, color) VALUES
('quick', '#10B981'),
('vegetarian', '#F59E0B'),
('comfort food', '#EF4444'),
('healthy', '#22C55E'),
('dessert', '#EC4899'),
('breakfast', '#8B5CF6'),
('dinner', '#3B82F6'),
('lunch', '#06B6D4'),
('grilled', '#F97316'),
('baked', '#A855F7');

-- Insert sample user
INSERT INTO users (email, password_hash, name, is_admin) VALUES
('demo@recipeapp.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Demo User', true);

-- Insert sample recipes
INSERT INTO recipes (user_id, title, description, instructions, prep_time, cook_time, servings, difficulty, is_public) VALUES
(1, 'Spaghetti Carbonara', 'Classic Italian pasta dish with eggs, cheese, and pancetta', 
'1. Cook spaghetti according to package directions. 2. Meanwhile, cook pancetta until crisp. 3. Beat eggs with parmesan cheese. 4. Drain pasta, reserve pasta water. 5. Mix hot pasta with pancetta. 6. Remove from heat, add egg mixture, toss quickly. 7. Add pasta water to achieve desired consistency.',
15, 20, 4, 'medium', true),
(1, 'Grilled Chicken Salad', 'Healthy grilled chicken with fresh vegetables', 
'1. Season chicken breast with salt, pepper, and herbs. 2. Grill chicken for 6-7 minutes per side. 3. Let chicken rest for 5 minutes. 4. Mix lettuce, tomatoes, cucumbers in a bowl. 5. Slice chicken and add to salad. 6. Dress with olive oil and lemon juice.',
10, 15, 2, 'easy', true);

-- Link ingredients to recipes
INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity, unit) VALUES
-- Spaghetti Carbonara ingredients
(1, 16, 400, 'grams'),    -- spaghetti
(1, 4, 200, 'ml'),       -- milk
(1, 3, 3, 'pieces'),     -- eggs
(1, 5, 100, 'grams'),    -- butter
(1, 6, 1, 'tsp'),        -- salt
(1, 7, 1, 'tsp'),        -- pepper
(1, 8, 2, 'cloves'),     -- garlic
(1, 18, 50, 'grams'),    -- parmesan cheese
-- Grilled Chicken Salad ingredients
(2, 14, 200, 'grams'),    -- chicken breast
(2, 20, 1, 'piece'),     -- bell peppers
(2, 11, 2, 'pieces'),     -- tomatoes
(2, 6, 1, 'tsp'),        -- salt
(2, 7, 1, 'tsp');        -- pepper

-- Link tags to recipes
INSERT INTO recipe_tags (recipe_id, tag_id) VALUES
(1, 1),   -- quick
(1, 3),   -- comfort food
(1, 8),   -- dinner
(2, 4),   -- healthy
(2, 7),   -- dinner
(2, 9);   -- grilled

-- Insert sample collections
INSERT INTO collections (user_id, name, description) VALUES
(1, 'Quick Meals', 'Recipes that take less than 30 minutes'),
(1, 'Family Favorites', 'Recipes everyone loves'),
(1, 'Healthy Options', 'Nutritious and delicious meals');

-- Add recipes to collections
INSERT INTO recipe_collections (recipe_id, collection_id) VALUES
(1, 2),   -- Spaghetti Carbonara -> Family Favorites
(2, 1),   -- Grilled Chicken Salad -> Quick Meals
(2, 3);   -- Grilled Chicken Salad -> Healthy Options