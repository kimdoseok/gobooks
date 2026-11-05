# Local Storage in Yew

## Direct Access
### in cargo.toml

[dependencies]
yew = { version = "0.22", features = ["csr"] }
serde = { version = "1.0", features = ["derive"] }
gloo-storage = "0.3" # Check for the latest version on crates.io

### in main.rs
```
use yew::prelude::*;
use gloo_storage::{LocalStorage, Storage};
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq)]
struct UserSession {
    username: String,
    token: String,
}

#[function_component(App)]
pub fn app() -> Html {
    // 1. Read from LocalStorage on component initialization
    let session_state = use_state(|| {
        LocalStorage::get::<UserSession>("user_session").unwrap_or_else(|_| UserSession {
            username: "Guest".to_string(),
            token: "".to_string(),
        })
    });

    // 2. Write to LocalStorage on action
    let login = {
        let session_state = session_state.clone();
        Callback::from(move |_| {
            let new_session = UserSession {
                username: "Ferris".to_string(),
                token: "secret_w@sm_token".to_string(),
            };
            
            // Save to browser storage
            if LocalStorage::set("user_session", &new_session).is_ok() {
                session_state.set(new_session);
            }
        })
    };

    // 3. Clear item from LocalStorage
    let logout = {
        let session_state = session_state.clone();
        Callback::from(move |_| {
            LocalStorage::delete("user_session");
            session_state.set(UserSession {
                username: "Guest".to_string(),
                token: "".to_string(),
            });
        })
    };

    html! {
        <div>
            <h1>{ format!("Welcome, {}", session_state.username) }</h1>
            { if session_state.token.is_empty() {
                html! { <button onclick={login}>{ "Log In" }</button> }
            } else {
                html! { <button onclick={logout}>{ "Log Out" }</button> }
            }}
        </div>
    }
}

```

## Reactive state

### In cargo.toml
```
[dependencies]
yew = { version = "0.22", features = ["csr"] }
yew-hooks = "0.6"

```

### in main.rs
```
use yew::prelude::*;
use yew_hooks::prelude::*;

#[function_component(LocalStorageHookExample)]
pub fn local_storage_hook_example() -> Html {
    // Initializes state with the value from localStorage key "username_key".
    // Defaults to "DefaultName" if the key doesn't exist yet.
    let storage = use_local_storage::<String>("username_key".to_string());

    let on_input_change = {
        let storage = storage.clone();
        Callback::from(move |e: InputEvent| {
            let input: web_sys::HtmlInputElement = e.target_unchecked_into();
            // This updates both the Yew state and the browser's localStorage automatically
            storage.set(input.value());
        })
    };

    html! {
        <div>
            <h3>{ "Stored Value: " }{ &*storage }</h3>
            <input type="text" oninput={on_input_change} value={(*storage).clone()} />
            <button onclick={move |_| storage.delete()}>{ "Reset Storage" }</button>
        </div>
    }
}

```

## Note for SSR
* calls to LocalStorage will panic 
* wrap your storage interactions inside a use_effect_with hook that guarantees to run in client 
