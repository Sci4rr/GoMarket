import React, { useEffect, useState } from 'react';

interface IProduct {
  id: number;
  name: string;
  category: string;
  price: number;
}

interface ICategory {
  name: string;
}

const ProductList: React.FC = () => {
  const [allProducts, setAllProducts] = useState<IProduct[]>([]);
  const [categoryList, setCategoryList] = useState<ICategory[]>([]);
  const [displayedProducts, setDisplayedProducts] = useState<IProduct[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [currentCategory, setCurrentCategory] = useState('All');
  const [sortCriteria, setSortCriteria] = useState('');

  useEffect(() => {
    async function fetchData() {
      const { products, categories } = await fetchAllData();
      setAllProducts(products);
      setDisplayedProducts(products);
      const allCategories = [{ name: 'All' }, ...categories];
      setCategoryList(allCategories);
    }

    fetchData();
  }, []);

  useEffect(() => {
    let filtered = [...allProducts];

    if (currentCategory !== 'All') {
      filtered = filtered.filter(product => product.category === currentCategory);
    }

    if (searchQuery) {
      filtered = filtered.filter(product =>
        product.name.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }

    switch (sortCriteria) {
      case 'priceLowHigh':
        filtered.sort((a, b) => a.price - b.price);
        break;
      case 'priceHighLow':
        filtered.sort((a, b) => b.price - a.price);
        break;
      default:
        break;
    }

    setDisplayedProducts(filtered);
  }, [searchQuery, currentCategory, sortCriteria, allProducts]);

  const handleSearchInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchQuery(event.target.value);
  };

  const handleCategorySelectionChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setCurrentCategory(event.target.value);
  };

  const handleSortingOptionChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setSortCriteria(event.target.value);
  };

  return (
    <div>
      <input type="text" placeholder="Search products..." onChange={handleSearchInputChange} value={searchQuery} />
      <select onChange={handleCategorySelectionChange} value={currentCategory}>
        {categoryList.map(category => (
          <option key={category.name} value={category.name}>
            {category.name}
          </option>
        ))}
      </select>
      <select onChange={handleSortingOptionChange} value={sortCriteria}>
        <option value="">Select Sorting Option</option>
        <option value="priceLowHigh">Price: Low to High</option>
        <option value="priceHighLow">Price: High to Low</option>
      </select>
      <ul>
        {displayedProducts.map(product => (
          <li key={product.id}>{`${product.name} - ${product.category} - $${product.price}`}</li>
        ))}
      </ul>
    </div>
  );
};

export default ProductList;

async function fetchAllData(): Promise<{ products: IProduct[], categories: ICategory[] }> {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve({
        products: [
          { id: 1, name: 'Laptop', category: 'Electronics', price: 1000 },
          { id: 2, name: 'Shirt', category: 'Clothing', price: 50 },
        ],
        categories: [
          { name: 'Electronics' },
          { name: 'Clothing' },
        ],
      });
    }, 1000);
  });
}