//query endpoint "/user_query?=user"
async function search(endpoint, query, abort) {
    const response = await fetch(`/${endpoint}?q=${query}`, {signal: abort});
    const input = await response.text();
    if (response.ok) {
        return ;
    } else {
        throw new Error(input);
    }
}