class App extends React.Component {

    state = {
        isLoggedIn: false,
        username: "",
        password: ""
    }

    componentDidMount() {
        this.checkAuth();
    }

    async checkAuth() {
        const res = await fetch("http://localhost:8080/CheckAuth");
        const resJson = await res.json();
        this.setState({isLoggedIn: resJson.isLoggedIn});
    }

    async handleSubmit() {

        const opts = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({username: this.state.username, password: this.state.password})
        }

        const res = await fetch("http://localhost:8080/Login", opts);
        const resJson = await res.json();
        this.setState({isLoggedIn: resJson.isLoggedIn})
    }

    render() {
        let content;
        
        if (this.state.isLoggedIn) {
            content = (
                <div>
                    <p>Hello, Rian</p>
                    <button onClick={async () => {
                        const res = await fetch("http://localhost:8080/Logout");
                        if (res.ok) {
                            this.setState({isLoggedIn: false});
                        }
                    }}>logout</button>
                </div>
            );
        } else {
            content = (
                <form onSubmit={e => {
                    e.preventDefault();
                    this.handleSubmit();
                }}>
                    <label>
                        username
                        <input type="text" value={this.state.username} onChange={e => {
                            this.setState({username: e.target.value})
                        }} />
                    </label>
                    <label>
                        password
                        <input type="text" value={this.state.password} onChange={e => {
                            this.setState({password: e.target.value})
                        }} />
                    </label>
                    <input type="submit" value="Submit" />
                </form>
            );
        }

        return (
            <div>
                <p>hello world from app</p>
                { content }
            </div>
        );
    }
}

const domContainer = document.getElementById('app');
const root = ReactDOM.createRoot(domContainer);
root.render(<App />);