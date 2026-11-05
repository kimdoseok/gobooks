# Single Level Router

## Cargo.toml
```
[dependencies]
yew = { version = "0.21", features = ["csr"] }
yew-router = "0.18"
```
## Router Enum
```
use yew_router::prelude::*;

#[derive(Clone, Routable, PartialEq)]
pub enum Route {
    #[at("/")]
    Home,
    #[at("/invoices")]
    Invoices,
    #[at("/expenses")]
    Expenses,
    #[at("/settings")]
    Settings,
    #[not_found]
    #[at("/404")]
    NotFound,
}
```

## Switch function
```
use yew::prelude::*;
use yew_router::prelude::*;

fn switch(routes: Route) -> Html {
    match routes {
        Route::Home => html! { <h1>{ "Dashboard Home" }</h1> },
        Route::Invoices => html! { <h1>{ "Invoices Overview" }</h1> },
        Route::Expenses => html! { <h1>{ "Expenses Tracker" }</h1> },
        Route::Settings => html! { <h1>{ "Settings" }</h1> },
        Route::NotFound => html! { <h1>{ "404 - Page Not Found" }</h1> },
    }
}
```

## main component
```
#[function_component(App)]
pub fn app() -> Html {
    html! {
        <BrowserRouter>
            {/* Nav links stay visible on every page */}
            <nav class="navbar">
                <Link<Route> to={Route::Home}>{ "Home" }</Link<Route>>
                <Link<Route> to={Route::Invoices}>{ "Invoices" }</Link<Route>>
                <Link<Route> to={Route::Expenses}>{ "Expenses" }</Link<Route>>
            </nav>

            <main class="content">
                {/* Switch renders whichever page matches the current URL */}
                <Switch<Route> render={switch} />
            </main>
        </BrowserRouter>
    }
}
```

## Navigation
### Static Link
```
use yew_router::prelude::*;

html! {
    <Link<Route> to={Route::Invoices} classes="nav-link">
        { "Go to Invoices" }
    </Link<Route>>
}
```

### Programtic Navigation
```
use yew::prelude::*;
use yew_router::prelude::*;

#[function_component(CreateInvoiceButton)]
pub fn create_invoice_button() -> Html {
    let navigator = use_navigator().unwrap();

    let onclick = Callback::from(move |_| {
        // Perform logic here...
        navigator.push(&Route::Invoices);
    });

    html! {
        <button {onclick}>{ "Save & View Invoices" }</button>
    }
}
```

## Route parameter
```
#[derive(Clone, Routable, PartialEq)]
pub enum Route {
    #[at("/")]
    Home,
    
    // Captures /invoices/123 -> InvoiceDetail { id: 123 }
    #[at("/invoices/:id")]
    InvoiceDetail { id: u32 },
    
    #[not_found]
    #[at("/404")]
    NotFound,
}

fn switch(routes: Route) -> Html {
    match routes {
        Route::Home => html! { <h1>{ "Home" }</h1> },
        Route::InvoiceDetail { id } => html! { 
            <h1>{ format!("Viewing Invoice #{}", id) }</h1> 
        },
        Route::NotFound => html! { <h1>{ "404" }</h1> },
    }
}
```

# Multi Level Router
## Parent and child Enum
```
use yew_router::prelude::*;

// 1. Top-Level Routes
#[derive(Clone, Routable, PartialEq, Debug)]
pub enum MainRoute {
    #[at("/")]
    Landing,
    
    // Catch-all wildcards (/*) for sub-routes inside /app
    #[at("/app/*")]
    AppRoot,

    #[not_found]
    #[at("/404")]
    NotFound,
}

// 2. Nested Sub-Routes (Internal Dashboard)
#[derive(Clone, Routable, PartialEq, Debug)]
pub enum AppSubRoute {
    #[at("/app/dashboard")]
    Dashboard,
    #[at("/app/invoices")]
    Invoices,
    #[at("/app/expenses")]
    Expenses,
    #[not_found]
    #[at("/app/404")]
    NotFound,
}
```

## Sub route switch
```
use yew::prelude::*;
use yew_router::prelude::*;

// Switch function specifically for the sub-routes
fn switch_app_sub_routes(route: AppSubRoute) -> Html {
    match route {
        AppSubRoute::Dashboard => html! { <h2>{ "Overview Dashboard" }</h2> },
        AppSubRoute::Invoices => html! { <h2>{ "Invoices & Billing" }</h2> },
        AppSubRoute::Expenses => html! { <h2>{ "Expense Management" }</h2> },
        AppSubRoute::NotFound => html! { <h2>{ "Sub-page not found" }</h2> },
    }
}

#[function_component(AppLayout)]
pub fn app_layout() -> Html {
    html! {
        <div class="dashboard-shell" style="display: flex; min-height: 100vh;">
            <!-- PERSISTENT SIDEBAR -->
            <aside style="width: 200px; background: #0f172a; color: white; padding: 1rem;">
                <h3>{ "rbooks" }</h3>
                <nav style="display: flex; flex-direction: column; gap: 0.5rem;">
                    <Link<AppSubRoute> to={AppSubRoute::Dashboard}>{ "Overview" }</Link<AppSubRoute>>
                    <Link<AppSubRoute> to={AppSubRoute::Invoices}>{ "Invoices" }</Link<AppSubRoute>>
                    <Link<AppSubRoute> to={AppSubRoute::Expenses}>{ "Expenses" }</Link<AppSubRoute>>
                </nav>
            </aside>

            <!-- NESTED CONTENT AREA -->
            <main style="flex: 1; padding: 2rem;">
                {/* Embedded Switch for Sub-Routes */}
                <Switch<AppSubRoute> render={switch_app_sub_routes} />
            </main>
        </div>
    }
}
```
## Connect main and sub routes
```
use yew::prelude::*;
use yew_router::prelude::*;

// Main Switch function for top-level pages
fn switch_main(route: MainRoute) -> Html {
    match route {
        MainRoute::Landing => html! {
            <div>
                <h1>{ "Welcome to rbooks" }</h1>
                <Link<AppSubRoute> to={AppSubRoute::Dashboard}>{ "Go to App" }</Link<AppSubRoute>>
            </div>
        },
        // Render the nested layout container for any /app/* route
        MainRoute::AppRoot => html! { <AppLayout /> },
        MainRoute::NotFound => html! { <h1>{ "404 - Page Not Found" }</h1> },
    }
}

#[function_component(App)]
pub fn app() -> Html {
    html! {
        <BrowserRouter>
            <Switch<MainRoute> render={switch_main} />
        </BrowserRouter>
    }
}
```
## Key points
1.	Wildcard /* Suffix: Ensure the parent route covering sub-pages uses /* (e.g., #[at("/app/*")]). Without the wildcard, yew-router won't pass control down to the inner .
2.	Full Path Matching in Sub-Routes: In AppSubRoute, paths must specify the full route string (e.g., #[at("/app/invoices")]).
3.	Cross-Enum Navigation: You can pass AppSubRoute variants directly into > anywhere inside , even from top-level landing pages.


# Dropdown Navigation Menu
## Hover dropdown
```
<!DOCTYPE html>
<html lang="en">
<head>
<style>
/* Container needed to position the dropdown content */
.dropdown {
  position: relative;
  display: inline-block;
}

/* Style the button that triggers the dropdown */
.dropbtn {
  background-color: #3498db;
  color: white;
  padding: 12px 24px;
  font-size: 16px;
  border: none;
  cursor: pointer;
  border-radius: 4px;
}

/* Dropdown Content (Hidden by Default) */
.dropdown-content {
  display: none;
  position: absolute;
  background-color: #f1f1f1;
  min-width: 160px;
  box-shadow: 0px 8px 16px rgba(0,0,0,0.2);
  z-index: 1;
  border-radius: 4px;
}

/* Links inside the dropdown */
.dropdown-content a {
  color: black;
  padding: 12px 16px;
  text-decoration: none;
  display: block;
}

/* Change color of dropdown links on hover */
.dropdown-content a:hover {
  background-color: #ddd;
}

/* Show the dropdown menu when hovering over the container */
.dropdown:hover .dropdown-content {
  display: block;
}

/* Change the button color when the dropdown is shown */
.dropdown:hover .dropbtn {
  background-color: #2980b9;
}
</style>
</head>
<body>

<div class="dropdown">
  <button class="dropbtn">Dropdown Menu</button>
  <div class="dropdown-content">
    <a href="#link1">Link Option 1</a>
    <a href="#link2">Link Option 2</a>
    <a href="#link3">Link Option 3</a>
  </div>
</div>

</body>
</html>

```
## Click to Open Dropdown
```
<!DOCTYPE html>
<html lang="en">
<head>
<style>
.dropdown {
  position: relative;
  display: inline-block;
}
.dropbtn {
  background-color: #2ecc71;
  color: white;
  padding: 12px 24px;
  font-size: 16px;
  border: none;
  cursor: pointer;
  border-radius: 4px;
}
.dropdown-content {
  display: none; /* Hidden by default */
  position: absolute;
  background-color: #f9f9f9;
  min-width: 160px;
  box-shadow: 0px 8px 16px rgba(0,0,0,0.2);
  z-index: 1;
}
.dropdown-content a {
  color: black;
  padding: 12px 16px;
  text-decoration: none;
  display: block;
}
.dropdown-content a:hover { background-color: #f1f1f1; }

/* This class is toggled by JavaScript */
.show { display: block; }
</style>
</head>
<body>

<div class="dropdown">
  <button id="myBtn" class="dropbtn">Click Me</button>
  <div id="myDropdown" class="dropdown-content">
    <a href="#link1">Option 1</a>
    <a href="#link2">Option 2</a>
    <a href="#link3">Option 3</a>
  </div>
</div>

<script>
const dropdownBtn = document.getElementById("myBtn");
const dropdownContent = document.getElementById("myDropdown");

// Toggle the menu when clicking the button
dropdownBtn.addEventListener("click", function(e) {
  dropdownContent.classList.toggle("show");
  e.stopPropagation(); // Prevents immediate closing from the window click listener
});

// Close the menu if the user clicks outside of it
window.addEventListener("click", function() {
  if (dropdownContent.classList.contains('show')) {
    dropdownContent.classList.remove('show');
  }
});
</script>

</body>
</html>

```

## Multi Level Dropdown(CSS only)
```
<!DOCTYPE html>
<html lang="en">
<head>
<style>
/* Main Navigation Bar */
ul.nav {
  list-style-type: none;
  margin: 0;
  padding: 0;
  background-color: #333;
}
ul.nav > li {
  float: left;
  position: relative;
}
ul.nav a {
  display: block;
  color: white;
  padding: 14px 16px;
  text-decoration: none;
}
ul.nav li:hover { background-color: #111; }

/* First Level Dropdown */
ul.dropdown {
  display: none;
  position: absolute;
  top: 100%;
  left: 0;
  background-color: #333;
  list-style: none;
  padding: 0;
  margin: 0;
  min-width: 160px;
}
ul.nav li:hover > ul.dropdown { display: block; }

/* Second Level Dropdown (Flyout) */
ul.submenu {
  display: none;
  position: absolute;
  top: 0;
  left: 100%; /* Positions it to the right of the parent menu */
  background-color: #444;
  list-style: none;
  padding: 0;
  margin: 0;
  min-width: 160px;
}
ul.dropdown li:position-relative { position: relative; } /* Context for submenu */
ul.dropdown li:hover > ul.submenu { display: block; }
</style>
</head>
<body>

<ul class="nav">
  <li><a href="#home">Home</a></li>
  <li style="position: relative;">
    <a href="#services">Services ▾</a>
    <!-- Level 1 Dropdown -->
    <ul class="dropdown">
      <li><a href="#design">Web Design</a></li>
      <li style="position: relative;">
        <a href="#dev">Development ▸</a>
        <!-- Level 2 Submenu -->
        <ul class="submenu">
          <li><a href="#frontend">Frontend</a></li>
          <li><a href="#backend">Backend</a></li>
        </ul>
      </li>
      <li><a href="#seo">SEO Marketing</a></li>
    </ul>
  </li>
  <li><a href="#contact">Contact</a></li>
</ul>

</body>
</html>

```