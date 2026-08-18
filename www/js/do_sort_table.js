document.querySelectorAll('#contents table:has(td)').forEach(table => {
  const tbody = table.querySelector('tbody') || table;
  if (tbody == null) { return; }
  const thead = table.querySelector('thead') || table.querySelector('tfoot') || table;
  if (thead == null) { return; }

  const rows = Array.from(tbody.querySelectorAll('tr:has(td)'));
  const headers = thead.querySelectorAll('tr:first-of-type th');
  
  rows.forEach((row, i) => row.dataset.originIndex = i);

  headers.forEach(th => {
    th.style.cursor = 'pointer';
    th.dataset.sortState = "0"; // 0:origal, 1:asc, 2:des

    th.addEventListener('click', () => {
      const index = Array.from(th.parentNode.children).indexOf(th);
      const currentState = parseInt(th.dataset.sortState, 10);
      const nextState = (currentState + 1) % 3;

      headers.forEach(h => {
        if (h !== th) h.dataset.sortState = "0";
      });

      th.dataset.sortState = nextState.toString();

      rows.sort((a, b) => {
        if (nextState === 0) {
          return a.dataset.originIndex - b.dataset.originIndex;
        }

        const valA = a.children[index].textContent.trim();
        const valB = b.children[index].textContent.trim();
        const isNum = !isNaN(valA) && !isNaN(valB);

        if (nextState === 1) {
          return isNum ? valA - valB : valA.localeCompare(valB);
        } else {
          return isNum ? valB - valA : valB.localeCompare(valA);
        }
      });

      rows.forEach(row => tbody.appendChild(row));
    });
  });
});
