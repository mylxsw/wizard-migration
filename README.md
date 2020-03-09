# wizard-migration

Wizard 数据迁移工具，用于将 ShowDoc 的文档迁移到 Wizard。迁移只保留项目、文档、文档历史、附件。

## 使用说明

在执行迁移之前，请先初始化好 Wizard 项目，并且创建管理员账号。

### 1. 迁移数据库

    ./wizard-migration \
        -showdoc_db "/ShowDoc路径/Sqlite/showdoc.db.php" \
        -wizard_db "账号:密码@tcp(数据库IP:数据库端口)/wizard_migration" \
        -replace_url "http://showdoc.local.yunsom.space"  \
        -import_user_id=1 
        
> `import_user_id` 指定了导入到 wizard 之后的文档的创建人、修改人默认用户ID

### 2. 迁移用户上传资源（附件，图片）

    cp -r /Showdoc项目根目录/Public/Uploads /Wizard项目根目录/storage/app/public/showdoc