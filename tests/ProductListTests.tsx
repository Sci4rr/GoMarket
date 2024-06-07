import React, { useState, useEffect } from 'react';
import axios from 'axios';

const ProductList = () => {
  const [products, setProducts] = useState([]);
  const [filteredProducts, setFilteredProducts] = useState([]);
  const [categories, setCategories] = useState(['All Categories']); // New State for categories
  const [selectedCategory, setSelectedCategory] = useState('');
  const [loading, setLoading] = useState(false); // New State for loading
  const [error, setError] = useState(''); // New State for error handling

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const result = await axios(`${process.env.REACT_APP_API_URL}/products`);
        const uniqueCategories = ['All Categories', ...new Set(result.data.map((product) => product.category))];
        setProducts(result.data);
        setFilteredProducts(result.data); // Initially set filteredProducts to all products
        setCategories(uniqueCategories);
      } catch (error) {
        setError('Failed to fetch products');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    if (selectedCategory === '') {
      setFilteredProducts(products);
    } else {
      setFilteredProducts(products.filter(product => product.category === selectedCategory));
    }
  }, [selectedCategory, products]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>{error}</p>;

  return (
    <div>
      <select data-testid="category-select" onChange={(e) => setSelectedCategory(e.target.value)}>
        {categories.map((category, index) => (
          <option key={index} value={category === 'All Categories' ? '' : category}>{category}</option>
        ))}
      </select>
      <div>
        {filteredProducts.map(product => (
          <div key={product.id} data-testid="product-item">{product.name} - ${product.dynamicPrice()}</div>
        ))}
      </div>
    </div>
  );
};

export default ProductList;