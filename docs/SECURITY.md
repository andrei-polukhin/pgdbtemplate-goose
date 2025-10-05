# Security Policy

## 🔒 Security Overview

The `pgdbtemplate-goose` adapter is designed with security as a first-class concern.
This document outlines our security practices, vulnerability disclosure process,
and security considerations for users of this library.

## 🚨 Reporting Security Vulnerabilities

If you discover a security vulnerability in `pgdbtemplate-goose`,
please help us by reporting it responsibly.

### 📞 Contact Information

**Please DO NOT report security vulnerabilities through public GitHub issues.**

Instead, please report security vulnerabilities using GitHub's
private vulnerability reporting:

- **GitHub Security Advisories**: [Report a vulnerability](https://github.com/andrei-polukhin/pgdbtemplate-goose/security/advisories/new)
- **Benefits**: Private, secure, and tracked through GitHub's security features

### 📋 Disclosure Process

1. **Report**: Submit a vulnerability report via
  [GitHub Security Advisories](https://github.com/andrei-polukhin/pgdbtemplate-goose/security/advisories/new)
2. **Acknowledgment**: You will receive an acknowledgment within 48 hours
3. **Investigation**: We will investigate and provide regular updates (at least weekly)
4. **Fix**: Once confirmed, we will work on a fix and coordinate disclosure
5. **Public Disclosure**: We will publish a security advisory once the fix is available

### 📝 What to Include in Your Report

Please include the following information in the description
of your vulnerability report:

- **Description**: A clear description of the vulnerability
- **Impact**: Potential impact and severity
- **Steps to Reproduce**: Detailed reproduction steps
- **Mitigation**: Any known workarounds or mitigations
- **Contact Information**: How we can reach you for follow-up

### 🏆 Recognition

We appreciate security researchers who help keep our users safe.
With your permission, we will acknowledge your contribution in our
security advisories and CONTRIBUTORS document.

## 🛡️ Security Considerations

### Migration File Security

**⚠️ Migration File Trust**

Goose executes SQL files from the migrations directory. Ensure:

- Migration files come from trusted sources
- Review all migration files before deployment
- Use version control to track migration changes
- Implement code review for all migrations

```go
// ✅ SECURE: Use versioned migrations from trusted repository.
runner := pgdbtemplategoose.NewMigrationRunner("./migrations")
```

### SQL Injection in Migrations

**⚠️ User Input in Migrations**

Never include unvalidated user input in migration files:

```sql
-- ❌ DANGEROUS: Never do this.
CREATE TABLE users_$USER_INPUT (id SERIAL);

-- ✅ SAFE: Use static, reviewed SQL.
CREATE TABLE users (id SERIAL PRIMARY KEY);
```

### Connection Security

**🔐 Connection String Handling**

- Connection strings should never be logged or exposed
- Use environment variables or secure credential stores
- Avoid hardcoding credentials in source code

```go
// ✅ RECOMMENDED: Use environment variables.
connString := os.Getenv("DATABASE_URL")

// ❌ AVOID: Hardcoded credentials.
connString := "postgres://user:password@localhost/db"
```

**🔒 TLS Configuration**

Always configure TLS for production databases:

```go
// ✅ SECURE: Require TLS.
connString := "postgres://user:pass@host/db?sslmode=require"

// ✅ SECURE: Verify CA certificate.
connString := "postgres://user:pass@host/db?sslmode=verify-ca"
```

### Database Permissions

**🔑 Principle of Least Privilege**

The database user should have minimal required permissions:

**For template creation (initialization)**:
```sql
-- Minimal permissions for template creation
GRANT CREATE ON DATABASE postgres TO migration_user;
GRANT ALL PRIVILEGES ON DATABASE template_db TO migration_user;
```

**For test database creation**:
```sql
-- Only needs template usage
GRANT USAGE ON SCHEMA public TO test_user;
```

### Dependencies

**📦 Dependency Management**

This adapter depends on:
- `pgdbtemplate` - Core template management library
- `pgdbtemplate-pq` - PostgreSQL driver (lib/pq)
- `pressly/goose/v3` - Migration framework

Ensure all dependencies are kept up to date:

```bash
# Check for updates
go list -u -m all

# Update dependencies
go get -u ./...
go mod tidy
```

**🔍 Vulnerability Scanning**

Regularly scan for known vulnerabilities:

```bash
# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# Scan for vulnerabilities
govulncheck ./...
```

## 🔐 Best Practices

### 1. Secure Migration Storage

```bash
# ✅ Good: Migrations in version control
./migrations/
  001_create_users.sql
  002_add_email_index.sql

# ✅ Good: Protected directory
chmod 755 migrations/
chmod 644 migrations/*.sql
```

### 2. Review Process

- Require code review for all migrations
- Test migrations in staging environment first
- Use CI/CD to validate migrations
- Document breaking changes

### 3. Backup Strategy

```go
// Always test with backups available
tm.Initialize(ctx) // Creates template

// Create backup before production deployment
// Use pg_dump or similar tools
```

### 4. Error Handling

```go
// ✅ Proper error handling
runner := pgdbtemplategoose.NewMigrationRunner("./migrations")
if err := runner.RunMigrations(ctx, conn); err != nil {
    log.Error("Migration failed", "error", err)
    // Implement rollback strategy
    return err
}
```

## 📚 Additional Resources

- [pgdbtemplate Security Policy](https://github.com/andrei-polukhin/pgdbtemplate/blob/main/docs/SECURITY.md)
- [Goose Security Considerations](https://github.com/pressly/goose#security)
- [PostgreSQL Security Best Practices](https://www.postgresql.org/docs/current/security.html)
- [OWASP Database Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Database_Security_Cheat_Sheet.html)

## 📧 Contact

For security-related questions that are not vulnerabilities,
please open a public discussion on GitHub.

---

**Last Updated**: 2025-10-05
