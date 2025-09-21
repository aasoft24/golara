# GoLara - Laravel-style Go Framework

A beautiful Laravel-inspired web framework for Go language, built with Go 1.24.6.

## 🚀 Features

- **Laravel-like Syntax**: Familiar Laravel syntax and structure
- **Database ORM**: GORM-based Eloquent-like ORM
- **Blade-like Templates**: Powerful template engine with layouts
- **Artisan CLI**: Laravel-style command line interface
- **Routing System**: Expressive route definitions with middleware support
- **Session Management**: Multiple session drivers support
- **File Storage**: Unified filesystem abstraction
- **Migrations & Seeding**: Database version control system
- **Middleware Support**: HTTP middleware pipeline
- **Validation**: Request validation system

## 📋 Requirements

- Go 1.24.6 or higher
- MySQL 5.7+, PostgreSQL 9.5+, or SQLite 3.8.8+

## 🔧 Installation

```bash
# Clone the repository
git clone https://github.com/aasoft24/golara.git
cd golara

# Install dependencies
go mod download

# Setup environment
cp .env.example .env

# Generate application key
go run artisan.go key:generate

# Run migrations
go run artisan.go migrate

# Seed the database (optional)
go run artisan.go db:seed

# Start development server
go run artisan.go serve
