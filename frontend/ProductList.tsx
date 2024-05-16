import React, { useEffect, useState } from 'react';

interface Product {
  id: number;
  name: string;
  category: string;
  price: number;
}

interface Category {
  name: string;
}

const ProductList: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [filteredProducts, setFilteredProducts] = useState<Product[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [sortOption, setSortOption] = useState('');

  useEffect(() => {
    const fetchedProducts: Product[] = [
      { id: 1, name: 'Laptop', category: 'Electronics', price: 1000 },
      { id: 2, name: 'Shirt', category: 'Clothing', price: 50 },
    ];
    setProducts(fetchedProducts);
    setFilteredProducts(fetchedProducts);

    const fetchedCategories: Category[] = [
      { name: 'All' },
      { name: 'Electronics' },
      { name: 'Clothing' },
    ];
    setCategories(fetchedCategories);
  }, []);

  useEffect(() => {
    let result = [...products];

    if (selectedCategory !== 'All') {
      result = result.filter(product => product.category === selectedCategory);
    }

    if (searchTerm) {
      result = result.filter(product =>
        product.name.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    switch (sortOption) {
      case 'priceLowHigh':
        result.sort((a, b) => a.price - b.price);
        break;
      case 'priceHighLow':
        result.sort((a, b) => b.price - a.price);
        break;
      default:
        break;
    }

    setFilteredProducts(result);
  }, [searchTerm, selectedCategory, sortOption, products]);

  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(event.target.value);
  };

  const handleCategoryChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setSelectedCategory(event.target.value);
  };

  const handleSortChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setSortOption(event.target.value);
  };

  return (
    <div>
      <input type="text" placeholder="Search products..." onChange={handleSearchChange} value={searchTerm} />
      <select onChange={handleCategoryChange} value={selectedCategory}>
        {categories.map(category => (
          <option key={category.name} value={category.name}>
            {category.name}
          </option>
        ))}
      </select>
      <select onChange={handleSortChange} value={sortOption}>
        <option value="">Select Sorting Option</option>
        <option value="priceLowHigh">Price: Low to High</option>
        <option value="priceHighLow">Price: High to Low</option>
      </select>
      <ul>
        {filteredProducts.map(product => (
          <li key={product.id}>{`${product.name} - ${product.category} - $${product.price}`}</li>
        ))}
      </ul>
    </div>
  );
};

export default ProductList;