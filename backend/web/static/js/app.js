// RecipeApp JavaScript utilities

// Modal functions
function showLoginModal() {
    document.getElementById('loginModal').classList.remove('hidden');
    document.getElementById('loginModal').classList.add('flex');
}

function hideLoginModal() {
    document.getElementById('loginModal').classList.add('hidden');
    document.getElementById('loginModal').classList.remove('flex');
}

function showRegisterModal() {
    document.getElementById('registerModal').classList.remove('hidden');
    document.getElementById('registerModal').classList.add('flex');
}

function hideRegisterModal() {
    document.getElementById('registerModal').classList.add('hidden');
    document.getElementById('registerModal').classList.remove('flex');
}

// Ingredient management for recipe form
let ingredientCount = 1;

function addIngredient() {
    const ingredientsList = document.getElementById('ingredients-list');
    const newIngredient = document.createElement('div');
    newIngredient.className = 'ingredient-item grid md:grid-cols-4 gap-4 mb-4';
    newIngredient.innerHTML = `
        <input type="text" name="ingredients[${ingredientCount}].name" placeholder="Ingredient name" required
               class="px-3 py-2 border rounded-lg focus:outline-none focus:border-blue-500">
        <input type="text" name="ingredients[${ingredientCount}].quantity" placeholder="Quantity" required
               class="px-3 py-2 border rounded-lg focus:outline-none focus:border-blue-500">
        <input type="text" name="ingredients[${ingredientCount}].unit" placeholder="Unit" required
               class="px-3 py-2 border rounded-lg focus:outline-none focus:border-blue-500">
        <button type="button" onclick="removeIngredient(this)" class="bg-red-500 text-white px-3 py-2 rounded hover:bg-red-600">Remove</button>
    `;
    ingredientsList.appendChild(newIngredient);
    ingredientCount++;
}

function removeIngredient(button) {
    const ingredientItem = button.closest('.ingredient-item');
    ingredientItem.remove();
}

// Instruction management for recipe form
let instructionCount = 1;

function addInstruction() {
    const instructionsList = document.getElementById('instructions-list');
    const newInstruction = document.createElement('div');
    newInstruction.className = 'instruction-item mb-4';
    newInstruction.innerHTML = `
        <div class="flex gap-4">
            <span class="font-semibold">${instructionCount + 1}.</span>
            <textarea name="instructions[${instructionCount}]" placeholder="Step ${instructionCount + 1}" rows="2" required
                      class="flex-1 px-3 py-2 border rounded-lg focus:outline-none focus:border-blue-500"></textarea>
            <button type="button" onclick="removeInstruction(this)" class="bg-red-500 text-white px-3 py-2 rounded hover:bg-red-600">Remove</button>
        </div>
    `;
    instructionsList.appendChild(newInstruction);
    instructionCount++;
}

function removeInstruction(button) {
    const instructionItem = button.closest('.instruction-item');
    instructionItem.remove();
    // Renumber remaining instructions
    const instructions = document.querySelectorAll('.instruction-item');
    instructions.forEach((item, index) => {
        const span = item.querySelector('span');
        const textarea = item.querySelector('textarea');
        span.textContent = `${index + 1}.`;
        textarea.placeholder = `Step ${index + 1}`;
    });
    instructionCount = instructions.length;
}

// Mobile menu toggle
function toggleMobileMenu() {
    const mobileMenu = document.querySelector('.mobile-menu');
    if (mobileMenu) {
        mobileMenu.classList.toggle('active');
    } else {
        // Create mobile menu if it doesn't exist
        const nav = document.querySelector('nav.hidden.md\\:flex');
        const mobileMenuCopy = nav.cloneNode(true);
        mobileMenuCopy.classList.remove('hidden', 'md:flex');
        mobileMenuCopy.classList.add('mobile-menu', 'fixed', 'top-16', 'left-0', 'w-full', 'bg-white', 'shadow-lg', 'flex-col', 'p-4', 'z-40');
        document.querySelector('header').appendChild(mobileMenuCopy);
    }
}

// Close modals when clicking outside
document.addEventListener('click', function(event) {
    const loginModal = document.getElementById('loginModal');
    const registerModal = document.getElementById('registerModal');
    
    if (event.target === loginModal) {
        hideLoginModal();
    }
    if (event.target === registerModal) {
        hideRegisterModal();
    }
});

// HTMX event handlers
document.addEventListener('htmx:afterRequest', function(event) {
    const target = event.target;
    
    // Handle successful authentication
    if (target.id === 'loginModal' || target.id === 'registerModal') {
        if (event.detail.successful) {
            // Redirect to reload the page with new auth state
            window.location.reload();
        }
    }
    
    // Handle recipe creation
    if (target.id === 'form-message') {
        if (event.detail.successful) {
            // Redirect to recipes page on successful creation
            setTimeout(() => {
                window.location.href = '/recipes';
            }, 1500);
        }
    }
});

// Form validation helpers
function validateForm(formData) {
    const errors = [];
    
    // Basic validation
    if (!formData.get('title')?.trim()) {
        errors.push('Recipe title is required');
    }
    if (!formData.get('description')?.trim()) {
        errors.push('Description is required');
    }
    if (!formData.get('difficulty')) {
        errors.push('Difficulty level is required');
    }
    
    return errors;
}

// Show flash messages
function showFlashMessage(message, type = 'info') {
    const flashDiv = document.createElement('div');
    flashDiv.className = `fixed top-20 right-4 p-4 rounded-lg shadow-lg z-50 ${
        type === 'error' ? 'bg-red-500 text-white' :
        type === 'success' ? 'bg-green-500 text-white' :
        type === 'warning' ? 'bg-yellow-500 text-white' :
        'bg-blue-500 text-white'
    }`;
    flashDiv.textContent = message;
    
    document.body.appendChild(flashDiv);
    
    setTimeout(() => {
        flashDiv.remove();
    }, 5000);
}