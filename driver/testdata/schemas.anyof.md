
# Sample (1.0)
sample tests

## default


<a id="getTest"></a>
<details>
<summary><code>GET</code> <code><b>/test</b></code> <code><i>test endpoint</i></code></summary>


### Responses
> | status | description | schema |
> | --- | --- | --- |
> | 200 | OK | [Identifier](#componentsSchemasIdentifier) |

</details>



## Schemas

<a id="componentsSchemasIdentifier"></a>
<details>
<summary><code>Identifier</code></summary>

### Any Of
> | schema |
> | --- |
> | integer(int64) |
> | string |
> | [Uuid](#componentsSchemasUuid) |

</details>

<a id="componentsSchemasUuid"></a>
<details>
<summary><code>Uuid</code></summary>

</details>



