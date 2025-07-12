-- Sample Data for Sales Track Database
-- This file contains realistic sample data for development and testing
-- Run this after creating the database schema

-- Clear existing data (for development/testing only)
DELETE FROM sales_records;

-- Reset auto-increment counter
DELETE FROM sqlite_sequence WHERE name='sales_records';

-- ============================================================================
-- SAMPLE SALES RECORDS
-- ============================================================================
-- Realistic sales data spanning multiple months, stores, and vendors
-- Includes various price points and scenarios for comprehensive testing

INSERT INTO sales_records (store, vendor, date, description, sale_price, commission, remaining) VALUES

-- January 2024 Sales
('Downtown Store', 'Electronics Plus', '2024-01-05', 'Samsung 55" Smart TV', 899.99, 89.99, 810.00),
('Mall Location', 'Home & Garden Co', '2024-01-05', 'Outdoor Patio Set', 1299.00, 129.90, 1169.10),
('Downtown Store', 'Fashion Forward', '2024-01-08', 'Winter Coat Collection', 245.50, 24.55, 220.95),
('Westside Branch', 'Electronics Plus', '2024-01-10', 'Apple iPad Pro', 1099.00, 109.90, 989.10),
('Mall Location', 'Sports Authority', '2024-01-12', 'Exercise Equipment Bundle', 675.25, 67.53, 607.72),
('Downtown Store', 'Kitchen Essentials', '2024-01-15', 'Professional Blender Set', 189.99, 19.00, 170.99),
('Eastside Outlet', 'Fashion Forward', '2024-01-18', 'Designer Handbag', 425.00, 42.50, 382.50),
('Westside Branch', 'Home & Garden Co', '2024-01-20', 'Indoor Plant Collection', 156.75, 15.68, 141.07),
('Mall Location', 'Electronics Plus', '2024-01-22', 'Wireless Headphones', 299.99, 30.00, 269.99),
('Downtown Store', 'Sports Authority', '2024-01-25', 'Running Shoes Premium', 179.95, 18.00, 161.95),

-- February 2024 Sales
('Eastside Outlet', 'Kitchen Essentials', '2024-02-02', 'Stainless Steel Cookware', 389.99, 39.00, 350.99),
('Westside Branch', 'Fashion Forward', '2024-02-05', 'Spring Dress Collection', 125.50, 12.55, 112.95),
('Mall Location', 'Home & Garden Co', '2024-02-08', 'Bedroom Furniture Set', 1899.00, 189.90, 1709.10),
('Downtown Store', 'Electronics Plus', '2024-02-10', 'Gaming Console Bundle', 549.99, 55.00, 494.99),
('Eastside Outlet', 'Sports Authority', '2024-02-12', 'Fitness Tracker Watch', 249.95, 25.00, 224.95),
('Westside Branch', 'Kitchen Essentials', '2024-02-15', 'Coffee Machine Deluxe', 299.00, 29.90, 269.10),
('Mall Location', 'Fashion Forward', '2024-02-18', 'Casual Wear Bundle', 89.99, 9.00, 80.99),
('Downtown Store', 'Home & Garden Co', '2024-02-20', 'Garden Tool Set', 145.75, 14.58, 131.17),
('Eastside Outlet', 'Electronics Plus', '2024-02-22', 'Smartphone Accessories', 79.95, 8.00, 71.95),
('Westside Branch', 'Sports Authority', '2024-02-25', 'Yoga Equipment Kit', 134.50, 13.45, 121.05),

-- March 2024 Sales
('Mall Location', 'Kitchen Essentials', '2024-03-01', 'Dining Table Set', 799.99, 80.00, 719.99),
('Downtown Store', 'Fashion Forward', '2024-03-05', 'Summer Collection Preview', 199.95, 20.00, 179.95),
('Eastside Outlet', 'Home & Garden Co', '2024-03-08', 'Outdoor Lighting System', 345.00, 34.50, 310.50),
('Westside Branch', 'Electronics Plus', '2024-03-10', 'Laptop Computer', 1299.99, 130.00, 1169.99),
('Mall Location', 'Sports Authority', '2024-03-12', 'Basketball Equipment', 189.75, 19.00, 170.75),
('Downtown Store', 'Kitchen Essentials', '2024-03-15', 'Food Processor Pro', 229.99, 23.00, 206.99),
('Eastside Outlet', 'Fashion Forward', '2024-03-18', 'Accessories Collection', 67.50, 6.75, 60.75),
('Westside Branch', 'Home & Garden Co', '2024-03-20', 'Living Room Decor', 456.25, 45.63, 410.62),
('Mall Location', 'Electronics Plus', '2024-03-22', 'Smart Home Hub', 199.00, 19.90, 179.10),
('Downtown Store', 'Sports Authority', '2024-03-25', 'Tennis Racket Set', 289.95, 29.00, 260.95),

-- April 2024 Sales
('Eastside Outlet', 'Kitchen Essentials', '2024-04-02', 'Microwave Oven', 189.99, 19.00, 170.99),
('Westside Branch', 'Fashion Forward', '2024-04-05', 'Business Attire', 345.00, 34.50, 310.50),
('Mall Location', 'Home & Garden Co', '2024-04-08', 'Bathroom Renovation Kit', 1156.75, 115.68, 1041.07),
('Downtown Store', 'Electronics Plus', '2024-04-10', 'Digital Camera', 699.99, 70.00, 629.99),
('Eastside Outlet', 'Sports Authority', '2024-04-12', 'Golf Club Set', 899.95, 90.00, 809.95),
('Westside Branch', 'Kitchen Essentials', '2024-04-15', 'Refrigerator Upgrade', 1599.00, 159.90, 1439.10),
('Mall Location', 'Fashion Forward', '2024-04-18', 'Shoe Collection', 156.50, 15.65, 140.85),
('Downtown Store', 'Home & Garden Co', '2024-04-20', 'Patio Furniture', 789.99, 79.00, 710.99),
('Eastside Outlet', 'Electronics Plus', '2024-04-22', 'Tablet Device', 449.95, 45.00, 404.95),
('Westside Branch', 'Sports Authority', '2024-04-25', 'Swimming Pool Supplies', 234.75, 23.48, 211.27),

-- May 2024 Sales
('Mall Location', 'Kitchen Essentials', '2024-05-01', 'Dishwasher Installation', 899.00, 89.90, 809.10),
('Downtown Store', 'Fashion Forward', '2024-05-05', 'Mother''s Day Special', 125.99, 12.60, 113.39),
('Eastside Outlet', 'Home & Garden Co', '2024-05-08', 'Lawn Mower Electric', 567.50, 56.75, 510.75),
('Westside Branch', 'Electronics Plus', '2024-05-10', 'Sound System', 799.99, 80.00, 719.99),
('Mall Location', 'Sports Authority', '2024-05-12', 'Camping Gear Bundle', 389.95, 39.00, 350.95),
('Downtown Store', 'Kitchen Essentials', '2024-05-15', 'Espresso Machine', 445.00, 44.50, 400.50),
('Eastside Outlet', 'Fashion Forward', '2024-05-18', 'Graduation Outfits', 289.75, 29.00, 260.75),
('Westside Branch', 'Home & Garden Co', '2024-05-20', 'Garden Shed Kit', 1299.99, 130.00, 1169.99),
('Mall Location', 'Electronics Plus', '2024-05-22', 'Gaming Accessories', 156.50, 15.65, 140.85),
('Downtown Store', 'Sports Authority', '2024-05-25', 'Bicycle Premium', 689.95, 69.00, 620.95),

-- June 2024 Sales (Recent)
('Eastside Outlet', 'Kitchen Essentials', '2024-06-02', 'Summer Appliances', 234.99, 23.50, 211.49),
('Westside Branch', 'Fashion Forward', '2024-06-05', 'Summer Fashion Line', 178.50, 17.85, 160.65),
('Mall Location', 'Home & Garden Co', '2024-06-08', 'Pool Equipment', 899.00, 89.90, 809.10),
('Downtown Store', 'Electronics Plus', '2024-06-10', 'Air Conditioning Unit', 1299.99, 130.00, 1169.99),
('Eastside Outlet', 'Sports Authority', '2024-06-12', 'Summer Sports Gear', 245.75, 24.58, 221.17);

-- ============================================================================
-- VERIFY SAMPLE DATA
-- ============================================================================
-- Quick verification queries to ensure data was inserted correctly

-- Total records inserted
-- SELECT COUNT(*) as total_records FROM sales_records;

-- Date range of sample data
-- SELECT MIN(date) as earliest_date, MAX(date) as latest_date FROM sales_records;

-- Total sales value
-- SELECT SUM(sale_price) as total_sales, SUM(commission) as total_commission, SUM(remaining) as total_remaining FROM sales_records;

-- Records by store
-- SELECT store, COUNT(*) as record_count FROM sales_records GROUP BY store ORDER BY record_count DESC;

-- Records by vendor
-- SELECT vendor, COUNT(*) as record_count FROM sales_records GROUP BY vendor ORDER BY record_count DESC;
