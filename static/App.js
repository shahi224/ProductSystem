const API = "http://localhost:8080";

let currentPage = 1;
let totalPages = 1;
let currentSearch = "";

// Create Product 
function createProduct() {
    const name = document.getElementById("name").value;
    const image = document.getElementById("image").value;
    const user = document.getElementById("user").value;
    const hsn = document.getElementById("hsn").value;
    
    if (!name || !user) {
        alert("Product Name and Created User are required!");
        return;
    }

    const data = {
        product_name: name,
        product_image: image || "default-image.jpg",
        created_user: user,
        hsn_code: hsn || ""
    };
    
    fetch(API + "/products", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(res => res.json())
    .then(result => {
        const msgDiv = document.getElementById("msg");
        msgDiv.innerHTML = `<div class="success">
            ✅ ${result.message}<br>
            Product ID: ${result.product?.id}<br>
            SubVariant ID: ${result.subvariant_id}
        </div>`;
        
        document.getElementById("name").value = "";
        document.getElementById("image").value = "";
        document.getElementById("user").value = "";
        document.getElementById("hsn").value = "";
        
        if (result.subvariant_id) {
            document.getElementById("subid").value = result.subvariant_id;
        }
        
        console.log("Product created:", result);
    })
    .catch(err => {
        console.error(err);
        document.getElementById("msg").innerHTML = 
            `<div class="error">❌ Error: ${err.message}</div>`;
    });
}

function loadProducts(page = 1, search = "") {
    currentPage = page;
    currentSearch = search;
    
    let url = `${API}/products?page=${page}&limit=10`;
    if (search) {
        url += `&search=${encodeURIComponent(search)}`;
    }
    
    console.log("Fetching products from:", url);
    
    fetch(url)
        .then(response => {
            console.log("Response status:", response.status, response.statusText);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log("API Response:", data);
            
            const ul = document.getElementById("list");
            if (!ul) return;
            
            ul.innerHTML = "";
            
            if (!data.products || !Array.isArray(data.products) || data.products.length === 0) {
                ul.innerHTML = `<li style="text-align:center;padding:40px;">
                    <h3>📭 No products found</h3>
                    <p>${search ? 'Try a different search term or ' : ''}
                    <a href="/create">create your first product</a></p>
                </li>`;
                
                const paginationDiv = document.getElementById("pagination");
                if (paginationDiv) paginationDiv.innerHTML = "";
                updateStats(data.total_items || 0);
                return;
            }
            
            currentPage = data.page || 1;
            totalPages = data.total_pages || 1;
            
            data.products.forEach(p => {
                const createdDate = new Date(p.created_date || p.createdDate);
                const formattedDate = createdDate.toLocaleDateString();
                
                ul.innerHTML += `
                    <li style="margin:15px 0;padding:20px;border-radius:8px;background:#fff;box-shadow:0 2px 5px rgba(0,0,0,0.1);">
                        <div style="display:flex;justify-content:space-between;align-items:start;">
                            <div style="flex:1;">
                                <h3 style="margin:0 0 10px 0;color:#333;">${p.product_name || 'Unnamed Product'}</h3>
                                <div style="color:#666;line-height:1.6;">
                                    <div><strong>📦 Code:</strong> ${p.product_code || 'No Code'}</div>
                                    <div><strong>🏷️ HSN:</strong> ${p.hsn_code || 'N/A'}</div>
                                    <div><strong>👤 Created by:</strong> ${p.created_user || 'Unknown'}</div>
                                    <div><strong>📅 Created:</strong> ${formattedDate}</div>
                                    <div><strong>Status:</strong> ${p.active ? '<span style="color:green;">✓ Active</span>' : '<span style="color:red;">✗ Inactive</span>'}</div>
                                </div>
                            </div>
                            <div style="text-align:right;margin-left:20px;">
                                <div style="font-size:32px;font-weight:bold;color:#2196F3;">
                                    ${p.total_stock || 0}
                                </div>
                                <div style="font-size:14px;color:#999;margin-bottom:10px;">Total Stock</div>
                                <button onclick="viewSubVariants('${p.id}')" style="padding:5px 15px;font-size:12px;margin-right:5px;">
                                    📦 SubVariants
                                </button>
                            </div>
                        </div>
                    </li>`;
            });
            
            updatePaginationControls(data);
            updateStats(data.total_items || 0);
        })
        .catch(error => {
            console.error("Error loading products:", error);
            const ul = document.getElementById("list");
            if (ul) {
                ul.innerHTML = `
                    <li class="error">
                        ❌ Error loading products: ${error.message}<br>
                        <small>Check browser console for details</small>
                    </li>`;
            }
        });
}



// Load all subvariants 
function loadAllSubVariants() {
    fetch(API + "/products?page=1&limit=100")
        .then(res => res.json())
        .then(result => {
            const products = result.products || [];
            const select = document.getElementById("subvariantSelect");
            if (!select) return;
            
            select.innerHTML = '<option value="">Select SubVariant</option>';
            
            const promises = products.map(product => 
                fetch(API + `/products/${product.id}/subvariants`)
                    .then(res => res.json())
                    .then(data => ({
                        product,
                        subvariants: data.subvariants || []
                    }))
            );
            
            Promise.all(promises)
                .then(results => {
                    results.forEach(({product, subvariants}) => {
                        subvariants.forEach(sv => {
                            const option = document.createElement("option");
                            option.value = sv.id;
                            option.textContent = `${product.product_name} - ${sv.sku} (Stock: ${sv.stock})`;
                            select.appendChild(option);
                        });
                    });
                })
                .catch(err => {
                    console.error("Error:", err);
                    select.innerHTML = `<option value="">Error loading subvariants</option>`;
                });
        })
        .catch(err => {
            console.error("Error loading products:", err);
            const select = document.getElementById("subvariantSelect");
            if (select) select.innerHTML = '<option value="">Error loading</option>';
        });
}

// update pagination control
function updatePaginationControls(data) {
    const paginationDiv = document.getElementById("pagination");
    if (!paginationDiv) return;
    
    const page = data.page || 1;
    const totalPages = data.total_pages || 1;
    const hasPrev = data.has_prev || false;
    const hasNext = data.has_next || false;
    const totalItems = data.total_items || 0;
    
    let html = '<div style="display:flex;justify-content:center;align-items:center;gap:10px;margin:20px 0;">';
    
    // Previous button
    if (hasPrev) {
        html += `<button onclick="goToPage(${page - 1})" style="padding:8px 15px;">
                    ← Previous
                 </button>`;
    } else {
        html += `<button disabled style="padding:8px 15px;opacity:0.5;">← Previous</button>`;
    }
    
    // Page numbers
    const maxVisiblePages = 5;
    let startPage = Math.max(1, page - Math.floor(maxVisiblePages / 2));
    let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1);
    
    if (endPage - startPage + 1 < maxVisiblePages) {
        startPage = Math.max(1, endPage - maxVisiblePages + 1);
    }
    
    if (startPage > 1) {
        html += `<button onclick="goToPage(1)" style="padding:8px 12px;">1</button>`;
        if (startPage > 2) html += `<span style="padding:8px 12px;">...</span>`;
    }
    
    for (let i = startPage; i <= endPage; i++) {
        if (i === page) {
            html += `<button style="background:#4CAF50;color:white;padding:8px 12px;">${i}</button>`;
        } else {
            html += `<button onclick="goToPage(${i})" style="padding:8px 12px;">${i}</button>`;
        }
    }
    
    if (endPage < totalPages) {
        if (endPage < totalPages - 1) html += `<span style="padding:8px 12px;">...</span>`;
        html += `<button onclick="goToPage(${totalPages})" style="padding:8px 12px;">${totalPages}</button>`;
    }
    

    if (hasNext) {
        html += `<button onclick="goToPage(${page + 1})" style="padding:8px 15px;">
                    Next →
                 </button>`;
    } else {
        html += `<button disabled style="padding:8px 15px;opacity:0.5;">Next →</button>`;
    }
    
    html += `</div>`;
    

    html += `<div style="text-align:center;color:#666;margin-top:10px;">
                Page ${page} of ${totalPages} | 
                Total Products: ${totalItems}
             </div>`;
    
    paginationDiv.innerHTML = html;
}

// Go to a s[ecific page 
function goToPage(page) {
    if (page < 1 || page > totalPages) return;
    loadProducts(page, currentSearch);
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

// Update statistics
function updateStats(totalItems) {
    const statsDiv = document.getElementById("stats");
    if (statsDiv) {
        statsDiv.innerHTML = `
            <div style="display:flex;gap:20px;justify-content:center;margin:20px 0;">
                <div style="text-align:center;">
                    <div style="font-size:24px;font-weight:bold;color:#4CAF50;">${totalItems}</div>
                    <div style="font-size:12px;color:#666;">Total Products</div>
                </div>
                <div style="text-align:center;">
                    <div style="font-size:24px;font-weight:bold;color:#2196F3;">${currentPage}</div>
                    <div style="font-size:12px;color:#666;">Current Page</div>
                </div>
                <div style="text-align:center;">
                    <div style="font-size:24px;font-weight:bold;color:#FF9800;">${totalPages}</div>
                    <div style="font-size:12px;color:#666;">Total Pages</div>
                </div>
            </div>
        `;
    }
}

// Search products
function searchProducts() {
    const searchInput = document.getElementById("searchInput");
    if (!searchInput) return;
    
    const searchTerm = searchInput.value.trim();
    loadProducts(1, searchTerm);
}

// Clear search
function clearSearch() {
    const searchInput = document.getElementById("searchInput");
    if (searchInput) {
        searchInput.value = "";
    }
    loadProducts(1, "");
}

// View subvariants
function viewSubVariants(productId) {
    fetch(`${API}/products/${productId}/subvariants`)
        .then(res => res.json())
        .then(data => {
            alert(`Product: ${data.product?.product_name}\n\nSubVariants:\n${
                data.subvariants.map(sv => 
                    `SKU: ${sv.sku}\nStock: ${sv.stock}\nID: ${sv.id}\n---`
                ).join('\n')
            }`);
        })
        .catch(err => {
            alert("Error loading subvariants: " + err.message);
        });
}

// Stock Management functions (same as before)
function stockIn() {
    stockAction("/stock/in");
}

function stockOut() {
    stockAction("/stock/out");
}

function stockAction(url) {
    const data = {
        sub_variant_id: document.getElementById("subid").value,
        qty: Number(document.getElementById("qty").value)
    };
    
    fetch(API + url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
    .then(res => res.json())
    .then(result => {
        document.getElementById("msg").innerText = result.message || result.error;
    })
    .catch(err => {
        console.error(err);
        alert("Error: " + err.message);
    });
}

// Stock Report
function loadReport() {
    fetch(API + "/stock/report")
        .then(res => res.json())
        .then(data => {
            const tbody = document.getElementById("report");
            tbody.innerHTML = "";
            
            data.forEach(t => {
                tbody.innerHTML += `
                    <tr>
                        <td>${t.sub_variant_id}</td>
                        <td>${t.quantity}</td>
                        <td>${t.transaction_type}</td>
                        <td>${new Date(t.transaction_date).toLocaleString()}</td>
                    </tr>`;
            });
        });
}

document.addEventListener('DOMContentLoaded', function() {
    if (document.getElementById("list")) {
        loadProducts();
    }
    

    if (document.getElementById("report")) {
        loadReport();
    }
});

