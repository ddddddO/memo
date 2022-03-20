# Standard Module Structure
- https://www.terraform.io/language/modules/develop/structure
- TODO: module化をすすめる

# Tips
- `terraform state mv <old> <new>`
    - 予め、resource単位でリソースを作成済みで、後から複数のresourceをひとまとめにしてmodule化した。
    - moduleを呼び出すようmain.tfを記述しterraform applyをしたが、作成済みのリソースをdestroyしてから新しくリソースを作るようなplan結果がでた。
    - そうはしたくないため、tfstateの整合性を取りたいため、上記コマンドを実行して解消した。
        - ref: https://mozami.me/2018/05/12/terraform_state_mv.html