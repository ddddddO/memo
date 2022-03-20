# memo
private & public memo system


## Components
![1](/arc/arc.png)

### Back-end(Go)
- REST API
- Batch
    - `exposer`
        - Change memos from private to public.
    - `notified-cnt-incrementer`
        - Increase the number since the last update date of the memo.

### Front-end(Vue.js)

#### Login
![1](/arc/_screen/login_by_smartphone.png)

#### View memo list
![1](/arc/_screen/memos_by_smartphone.png)

#### View memo list(PC)
![1](/arc/_screen/memos_by_browser.png)

#### Edit memo
![1](/arc/_screen/memo_edit_by_smartphone.png?raw=true)
