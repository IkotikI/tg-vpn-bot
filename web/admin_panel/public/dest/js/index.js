let user;
user = {
    ID: 1,
    TelegramID: 1,
    TelegramName: "test",
    CreatedAt: new Date(),
    UpdatedAt: new Date(),
};
console.log(user);
/* ---- Darkmode ---- */
var btns = document.querySelectorAll(".darkmode-toggle-btn");
btns.forEach((btn) => {
    btn.addEventListener("click", () => {
        const html = document.documentElement;
        const transitionDelay = 200;
        const transitionClasses = ["transition-colors", ("duration-" + transitionDelay)];
        html.classList.add(...transitionClasses);
        document.documentElement.classList.toggle("dark");
        setTimeout(() => { html.classList.remove(...transitionClasses); }, transitionDelay);
    });
});
/* ---- Searchable Dropdown ---- */
const searchable_dropdowns = document.getElementsByClassName('searchable-dropdown');
Array.from(searchable_dropdowns).forEach((group) => {
    SearchableDropdown(group);
});
function SearchableDropdown(group) {
    const dropdownInput = group.getElementsByClassName('dropdown-input')[0];
    const dropdownButton = group.getElementsByClassName('dropdown-button')[0];
    const dropdownButtonContent = dropdownButton.getElementsByClassName('button-content')[0];
    const dropdownMenu = group.getElementsByClassName('dropdown-menu')[0];
    const searchInput = group.getElementsByClassName('search-input')[0];
    const dropdownMenuItems = group.getElementsByClassName('dropdown-menu-item');
    let isOpen = false; // Set to true to open the dropdown by default
    console.log(dropdownInput);
    if (dropdownButton == null || dropdownMenu == null || searchInput == null) {
        console.error("Dropdown inconsistent: some elements not found: (button, menu, input):", dropdownButton, dropdownMenu, searchInput);
        return;
    }
    // Function to toggle the dropdown state
    function toggleDropdown() {
        isOpen = !isOpen;
        dropdownMenu.classList.toggle('hidden', !isOpen);
    }
    // Set initial state
    // toggleDropdown();
    dropdownButton.addEventListener('click', () => {
        toggleDropdown();
        searchInput.focus();
    });
    // Add event listener to filter items based on input
    searchInput.addEventListener('input', () => {
        const searchTerm = searchInput.value.toLowerCase();
        const items = dropdownMenu.querySelectorAll('a');
        items.forEach((item) => {
            const text = item?.textContent?.toLowerCase();
            if (text?.includes(searchTerm)) {
                item.style.display = 'block';
            }
            else {
                item.style.display = 'none';
            }
        });
    });
    for (let i = 0; i < dropdownMenuItems.length; i++) {
        dropdownMenuItems[i].addEventListener('click', () => {
            toggleDropdown();
            dropdownButtonContent.innerHTML = dropdownMenuItems[i].innerHTML;
            dropdownInput.value = dropdownMenuItems[i].getAttribute('data-value') ?? "";
            console.log("dropdownInput.value", dropdownInput.value);
        });
    }
}
export {};
//# sourceMappingURL=index.js.map