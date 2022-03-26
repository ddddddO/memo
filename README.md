# memo
private & public memo system


## Components
![1](/doc/arc.svg)

### Back-end(Go)
- REST API
- Batch
    - `exposer`
        - Change memos from private to public.
    - `notified-cnt-incrementer`
        - Increase the number since the last update date of the memo.
        - TODO: remove or rename and refactor.

### Front-end(Vue.js)
- Manage private memos

#### Login
![1](/doc/screenshot/login_by_smartphone_new.png)

#### Memos/list
![1](/doc/screenshot/memos_by_smartphone_new.png)

#### Memos/list(PC)
![1](/doc/screenshot/memos_by_pc_new.png)

#### Memo/edit
![1](/doc/screenshot/memo_edit_by_smartphone_new_1.png)
![1](/doc/screenshot/memo_edit_by_smartphone_new_2.png)

#### Tags/list
![1](/doc/screenshot/tags_by_smartphone_new.png)

#### Tag/edit
![1](/doc/screenshot/tag_edit_by_smartphone_new.png)

### http://www.dododo.site/
- View public memos

#### Memos/list
![1](/doc/screenshot/public_memos_by_smartphone.png)

#### Memo/show
![1](/doc/screenshot/public_memo_by_smartphone.png)
